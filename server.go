package lambdabase

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-nacelle/nacelle"
)

type (
	LambdaServer struct {
		Logger   nacelle.Logger           `service:"logger"`
		Services nacelle.ServiceContainer `service:"services"`
		handler  Handler
		listener net.Listener
		server   *rpc.Server
		once     *sync.Once
	}

	Handler interface {
		nacelle.Initializer
		lambda.Handler
	}
)

func NewServer(handler Handler) *LambdaServer {
	return &LambdaServer{
		handler: handler,
		once:    &sync.Once{},
	}
}

func (s *LambdaServer) Init(config nacelle.Config) error {
	serverConfig := &Config{}
	if err := config.Load(serverConfig); err != nil {
		return err
	}

	if err := s.Services.Inject(s.handler); err != nil {
		return err
	}

	if err := s.handler.Init(config); err != nil {
		return err
	}

	server := rpc.NewServer()
	handler := lambda.NewFunction(s.handler)

	if err := server.Register(handler); err != nil {
		return fmt.Errorf("failed to register RPC (%s)", err.Error())
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", serverConfig.LambdaServerPort))
	if err != nil {
		return fmt.Errorf("failed to create listener (%s)", err.Error())
	}

	s.server = server
	s.listener = listener
	return nil
}

func (s *LambdaServer) Start() error {
	defer s.close()
	wg := sync.WaitGroup{}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok {
				if opErr.Err.Error() == "use of closed network connection" {
					break
				}
			}

			return err
		}

		wg.Add(1)

		go func() {
			defer wg.Done()
			s.server.ServeConn(conn)
		}()
	}

	s.Logger.Info("Draining lambda server")
	wg.Wait()
	return nil
}

func (s *LambdaServer) Stop() error {
	s.close()
	return nil
}

func (s *LambdaServer) close() {
	s.once.Do(func() {
		if s.listener == nil {
			return
		}

		s.Logger.Info("Closing lambda listener")
		s.listener.Close()
	})
}

package lambdabase

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-nacelle/nacelle"
)

type (
	Server struct {
		Logger   nacelle.Logger           `service:"logger"`
		Services nacelle.ServiceContainer `service:"services"`
		handler  Handler
		listener net.Listener
		server   *rpc.Server
		once     *sync.Once
	}

	// TODO - add health?

	Handler interface {
		nacelle.Initializer
		lambda.Handler
	}

	LambdaHandlerFunc func(ctx context.Context, payload []byte) ([]byte, error)
)

func (f LambdaHandlerFunc) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	return f(ctx, payload)
}

func NewServer(handler Handler) *Server {
	return &Server{
		handler: handler,
		once:    &sync.Once{},
	}
}

func (s *Server) Init(config nacelle.Config) error {
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

func (s *Server) Start() error {
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

func (s *Server) Stop() error {
	s.close()
	return nil
}

func (s *Server) close() {
	s.once.Do(func() {
		if s.listener == nil {
			return
		}

		s.Logger.Info("Closing lambda listener")
		s.listener.Close()
	})
}

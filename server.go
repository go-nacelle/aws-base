package lambdabase

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-nacelle/nacelle"
	"github.com/google/uuid"
)

type (
	Server struct {
		Logger      nacelle.Logger           `service:"logger"`
		Services    nacelle.ServiceContainer `service:"services"`
		Health      nacelle.Health           `service:"health"`
		handler     Handler
		listener    net.Listener
		server      *rpc.Server
		once        *sync.Once
		healthToken healthToken
	}

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
		handler:     handler,
		once:        &sync.Once{},
		healthToken: healthToken(uuid.New().String()),
	}
}

func (s *Server) Init(config nacelle.Config) error {
	if err := s.Health.AddReason(s.healthToken); err != nil {
		return err
	}

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

	listener, err := makeListener("", serverConfig.LambdaServerPort)
	if err != nil {
		return err
	}

	server := rpc.NewServer()

	if err := server.Register(lambda.NewFunction(s.handler)); err != nil {
		return fmt.Errorf("failed to register RPC (%s)", err.Error())
	}

	s.server = server
	s.listener = listener
	return nil
}

func (s *Server) Start() error {
	defer s.close()
	wg := sync.WaitGroup{}

	if err := s.Health.RemoveReason(s.healthToken); err != nil {
		return err
	}

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

func makeListener(host string, port int) (*net.TCPListener, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	return net.ListenTCP("tcp", addr)
}

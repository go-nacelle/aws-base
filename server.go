package lambdabase

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-nacelle/config/v3"
	"github.com/go-nacelle/nacelle/v2"
	"github.com/go-nacelle/process/v2"
	"github.com/go-nacelle/service/v2"
	"github.com/google/uuid"
)

type (
	Server struct {
		Config       *nacelle.Config           `service:"config"`
		Logger       nacelle.Logger            `service:"logger"`
		Services     *nacelle.ServiceContainer `service:"services"`
		Health       *nacelle.Health           `service:"health"`
		handler      Handler
		listener     net.Listener
		server       *rpc.Server
		once         *sync.Once
		healthToken  healthToken
		healthStatus *process.HealthComponentStatus
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

func (s *Server) Init(ctx context.Context) error {
	healthStatus, err := s.Health.Register(s.healthToken)
	if err != nil {
		return err
	}
	s.healthStatus = healthStatus

	serverConfig := &Config{}
	if err := config.LoadFromContext(ctx, serverConfig); err != nil {
		return err
	}

	if err := service.Inject(ctx, s.Services, s.handler); err != nil {
		return err
	}

	if err := s.handler.Init(ctx); err != nil {
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

func (s *Server) Run(ctx context.Context) error {
	defer s.close()
	wg := sync.WaitGroup{}

	s.healthStatus.Update(true)

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

func (s *Server) Stop(ctx context.Context) error {
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

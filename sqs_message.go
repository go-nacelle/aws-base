package lambdabase

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-nacelle/nacelle/v2"
)

type (
	SQSMessageHandler interface {
		Handle(ctx context.Context, message events.SQSMessage, logger nacelle.Logger) error
	}

	sqsMessageHandlerInitializer interface {
		nacelle.Initializer
		SQSMessageHandler
	}

	sqsMessageHandler struct {
		Logger   nacelle.Logger            `service:"logger"`
		Services *nacelle.ServiceContainer `service:"services"`
		handler  SQSMessageHandler
	}
)

func NewSQSRecordServer(handler SQSMessageHandler) *Server {
	return NewSQSEventServer(&sqsMessageHandler{
		handler: handler,
	})
}

func (s *sqsMessageHandler) Init(ctx context.Context) error {
	return doInit(ctx, s.Services, s.handler)
}

func (h *sqsMessageHandler) Handle(ctx context.Context, batch []events.SQSMessage, logger nacelle.Logger) error {
	for _, message := range batch {
		messageLogger := logger.WithFields(map[string]interface{}{
			"messageId": message.MessageId,
		})

		logger.Debug("Handling message")

		if err := h.handler.Handle(ctx, message, messageLogger); err != nil {
			return fmt.Errorf("failed to process SQS message %s (%s)", message.MessageId, err.Error())
		}
	}

	logger.Debug("SQS message handled successfully")
	return nil
}

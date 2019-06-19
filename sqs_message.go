package lambdabase

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-nacelle/nacelle"
)

type (
	SQSMessageHandler interface {
		Handle(ctx context.Context, message events.SQSMessage, logger nacelle.Logger) error
	}

	sqsMessageHandler struct {
		Logger   nacelle.Logger           `service:"logger"`
		Services nacelle.ServiceContainer `service:"services"`
		handler  SQSMessageHandler
	}
)

func NewSQSRecordServer(handler SQSMessageHandler) nacelle.Process {
	return NewSQSEventServer(&sqsMessageHandler{
		handler: handler,
	})
}

func (s *sqsMessageHandler) Init(config nacelle.Config) error {
	return doInit(config, s.Services, s.handler)
}

func (h *sqsMessageHandler) Handle(ctx context.Context, batch []events.SQSMessage, logger nacelle.Logger) error {
	for _, message := range batch {
		messageLogger := logger.WithFields(map[string]interface{}{
			"messageId": message.MessageId,
		})

		// TODO - log

		if err := h.handler.Handle(ctx, message, messageLogger); err != nil {
			// TODO - log
			return err
		}
	}

	// TODO - log
	return nil
}

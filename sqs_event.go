package lambdabase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-nacelle/nacelle/v2"
)

type (
	SQSEventHandler interface {
		Handle(ctx context.Context, batch []events.SQSMessage, logger nacelle.Logger) error
	}

	sqsEventHandlerInitializer interface {
		nacelle.Initializer
		SQSEventHandler
	}

	sqsEventHandler struct {
		Logger   nacelle.Logger            `service:"logger"`
		Services *nacelle.ServiceContainer `service:"services"`
		handler  SQSEventHandler
	}
)

func NewSQSEventServer(handler SQSEventHandler) *Server {
	return NewServer(&sqsEventHandler{
		handler: handler,
	})
}

func (h *sqsEventHandler) Init(ctx context.Context) error {
	return doInit(ctx, h.Services, h.handler)
}

func (h *sqsEventHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	event := &events.SQSEvent{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event (%s)", err.Error())
	}

	logger := h.Logger.WithFields(map[string]interface{}{
		"requestId": GetRequestID(ctx),
	})

	logger.Debug("Received %d SQS messages", len(event.Records))

	if err := h.handler.Handle(ctx, event.Records, logger); err != nil {
		return nil, fmt.Errorf("failed to process SQS event (%s)", err.Error())
	}

	logger.Debug("SQS event handled successfully")
	return nil, nil
}

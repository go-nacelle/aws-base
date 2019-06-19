package lambdabase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-nacelle/nacelle"
)

type (
	SQSEventHandler interface {
		Handle(ctx context.Context, batch []events.SQSMessage, logger nacelle.Logger) error
	}

	sqsEventHandler struct {
		Logger   nacelle.Logger           `service:"logger"`
		Services nacelle.ServiceContainer `service:"services"`
		handler  SQSEventHandler
	}
)

func NewSQSEventServer(handler SQSEventHandler) nacelle.Process {
	return NewServer(&sqsEventHandler{
		handler: handler,
	})
}

func (h *sqsEventHandler) Init(config nacelle.Config) error {
	return doInit(config, h.Services, h.handler)
}

func (h *sqsEventHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	event := &events.SQSEvent{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event (%s)", err.Error())
	}

	// TODO - log

	err := h.handler.Handle(ctx, event.Records, h.Logger.WithFields(map[string]interface{}{
		"requestId": getRequestID(ctx),
	}))

	if err != nil {
		// TODO -
		return nil, err
	}

	// TODO - log
	return nil, nil
}

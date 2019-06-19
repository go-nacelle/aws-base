package lambdabase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-nacelle/nacelle"
)

type (
	KinesisEventHandler interface {
		Handle(ctx context.Context, batch []events.KinesisEventRecord, logger nacelle.Logger) error
	}

	kinesisEventHandler struct {
		Logger   nacelle.Logger           `service:"logger"`
		Services nacelle.ServiceContainer `service:"services"`
		handler  KinesisEventHandler
	}
)

func NewKinesisEventServer(handler KinesisEventHandler) nacelle.Process {
	return NewServer(&kinesisEventHandler{
		handler: handler,
	})
}

func (h *kinesisEventHandler) Init(config nacelle.Config) error {
	return doInit(config, h.Services, h.handler)
}

func (h *kinesisEventHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	event := &events.KinesisEvent{}
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

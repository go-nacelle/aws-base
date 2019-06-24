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

	kinesisEventHandlerInitializer interface {
		nacelle.Initializer
		KinesisEventHandler
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

	logger := h.Logger.WithFields(map[string]interface{}{
		"requestId": GetRequestID(ctx),
	})

	logger.Debug("Received %d Kinesis records", len(event.Records))

	if err := h.handler.Handle(ctx, event.Records, logger); err != nil {
		return nil, fmt.Errorf("failed to process Kinesis event (%s)", err.Error())
	}

	logger.Debug("Kinesis event handled successfully")
	return nil, nil
}

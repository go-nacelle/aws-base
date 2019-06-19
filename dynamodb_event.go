package lambdabase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-nacelle/nacelle"
)

type (
	DynamoDBEventHandler interface {
		Handle(ctx context.Context, batch []events.DynamoDBEventRecord, logger nacelle.Logger) error
	}

	dynamoDBEventHandler struct {
		Logger   nacelle.Logger           `service:"logger"`
		Services nacelle.ServiceContainer `service:"services"`
		handler  DynamoDBEventHandler
	}
)

func NewDynamoDBEventServer(handler DynamoDBEventHandler) nacelle.Process {
	return NewServer(&dynamoDBEventHandler{
		handler: handler,
	})
}

func (h *dynamoDBEventHandler) Init(config nacelle.Config) error {
	return doInit(config, h.Services, h.handler)
}

func (h *dynamoDBEventHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	event := &events.DynamoDBEvent{}
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

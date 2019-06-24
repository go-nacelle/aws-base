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

	dynamoDBEventHandlerInitializer interface {
		nacelle.Initializer
		DynamoDBEventHandler
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

	logger := h.Logger.WithFields(map[string]interface{}{
		"requestId": GetRequestID(ctx),
	})

	logger.Debug("Received %d DynamoDB records", len(event.Records))

	if err := h.handler.Handle(ctx, event.Records, logger); err != nil {
		return nil, fmt.Errorf("failed to process DynamoDB event (%s)", err.Error())
	}

	logger.Debug("DynamoDB event handled successfully")
	return nil, nil
}

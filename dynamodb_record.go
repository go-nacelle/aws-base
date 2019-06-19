package lambdabase

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-nacelle/nacelle"
)

type (
	DynamoDBRecordHandler interface {
		Handle(ctx context.Context, record events.DynamoDBEventRecord, logger nacelle.Logger) error
	}

	dynamoDBRecordHandler struct {
		Logger   nacelle.Logger           `service:"logger"`
		Services nacelle.ServiceContainer `service:"services"`
		handler  DynamoDBRecordHandler
	}
)

func NewDynamoDBRecordServer(handler DynamoDBRecordHandler) nacelle.Process {
	return NewDynamoDBEventServer(&dynamoDBRecordHandler{
		handler: handler,
	})
}

func (s *dynamoDBRecordHandler) Init(config nacelle.Config) error {
	return doInit(config, s.Services, s.handler)
}

func (h *dynamoDBRecordHandler) Handle(ctx context.Context, records []events.DynamoDBEventRecord, logger nacelle.Logger) error {
	for _, record := range records {
		recordLogger := logger.WithFields(map[string]interface{}{
			"eventId": record.EventID,
		})

		// TODO - log

		if err := h.handler.Handle(ctx, record, recordLogger); err != nil {
			// TODO - log
			return err
		}
	}

	// TODO - log
	return nil
}

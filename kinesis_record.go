package lambdabase

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-nacelle/nacelle"
)

type (
	KinesisRecordHandler interface {
		Handle(ctx context.Context, record events.KinesisEventRecord, logger nacelle.Logger) error
	}

	kinesisRecordHandlerInitializer interface {
		nacelle.Initializer
		KinesisRecordHandler
	}

	kinesisRecordHandler struct {
		Logger   nacelle.Logger           `service:"logger"`
		Services nacelle.ServiceContainer `service:"services"`
		handler  KinesisRecordHandler
	}
)

func NewKinesisRecordServer(handler KinesisRecordHandler) nacelle.Process {
	return NewKinesisEventServer(&kinesisRecordHandler{
		handler: handler,
	})
}

func (s *kinesisRecordHandler) Init(config nacelle.Config) error {
	return doInit(config, s.Services, s.handler)
}

func (h *kinesisRecordHandler) Handle(ctx context.Context, records []events.KinesisEventRecord, logger nacelle.Logger) error {
	for _, record := range records {
		recordLogger := logger.WithFields(map[string]interface{}{
			"eventId": record.EventID,
		})

		logger.Debug("Handling record")

		if err := h.handler.Handle(ctx, record, recordLogger); err != nil {
			return fmt.Errorf("failed to process Kinesis record %s (%s)", record.EventID, err.Error())
		}
	}

	logger.Debug("Kinesis record handled successfully")
	return nil
}

package lambdabase

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-nacelle/nacelle"
)

type (
	KinesisRecordHandler interface {
		Handle(ctx context.Context, record events.KinesisEventRecord, logger nacelle.Logger) error
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

		// TODO - log

		if err := h.handler.Handle(ctx, record, recordLogger); err != nil {
			// TODO - log
			return err
		}
	}

	// TODO - log
	return nil
}

package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-nacelle/lambdabase"
	"github.com/go-nacelle/nacelle"
)

type RecordHandler struct{}

func (h *RecordHandler) Handle(ctx context.Context, record events.KinesisEventRecord, logger nacelle.Logger) error {
	logger.Info("Data: %s", string(record.Kinesis.Data))
	return nil
}

func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
	processes.RegisterProcess(lambdabase.NewKinesisRecordServer(&RecordHandler{}), nacelle.WithProcessName("kinesis-handler"))
	return nil
}

func main() {
	nacelle.NewBootstrapper("lambdabase-example", setup).BootAndExit()
}

package lambdabase

//go:generate go-mockgen -f github.com/go-nacelle/lambdabase -i dynamoDBEventHandlerInitializer -i dynamoDBRecordHandlerInitializer -o dynamodb_mock_test.go
//go:generate go-mockgen -f github.com/go-nacelle/lambdabase -i kinesisEventHandlerInitializer -i kinesisRecordHandlerInitializer -o kinesis_mock_test.go
//go:generate go-mockgen -f github.com/go-nacelle/lambdabase -i sqsEventHandlerInitializer -i sqsMessageHandlerInitializer -o sqs_mock_test.go

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

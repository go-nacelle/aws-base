package lambdabase

//go:generate go-mockgen -f github.com/go-nacelle/lambdabase -i dynamoDBEventHandlerInitializer -i dynamoDBRecordHandlerInitializer -o dynamodb_mock_test.go
//go:generate go-mockgen -f github.com/go-nacelle/lambdabase -i kinesisEventHandlerInitializer -i kinesisRecordHandlerInitializer -o kinesis_mock_test.go
//go:generate go-mockgen -f github.com/go-nacelle/lambdabase -i sqsEventHandlerInitializer -i sqsMessageHandlerInitializer -o sqs_mock_test.go

import (
	"testing"

	"github.com/aphistic/sweet"
	junit "github.com/aphistic/sweet-junit"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&DynamoDBSuite{})
		s.AddSuite(&KinesisSuite{})
		s.AddSuite(&ServerSuite{})
		s.AddSuite(&SQSSuite{})
	})
}

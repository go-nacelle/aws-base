package lambdabase

import (
	"context"
	"fmt"

	"github.com/aphistic/sweet"
	"github.com/aws/aws-lambda-go/events"
	. "github.com/efritz/go-mockgen/matchers"
	"github.com/go-nacelle/nacelle"
	. "github.com/onsi/gomega"
)

type DynamoDBSuite struct{}

var testDynamoDBPayload = `{
	"Records": [
		{
			"eventID": "ev1",
			"eventName": "INSERT",
			"dynamodb": {
				"NewImage": {
					"PK": {"S": "foo"},
					"SK": {"S": "bonk"}
				}
			}
		},
		{
			"eventID": "ev2",
			"eventName": "INSERT",
			"dynamodb": {
				"NewImage": {
					"PK": {"S": "bar"},
					"SK": {"S": "quux"}
				}
			}
		},
		{
			"eventID": "ev3",
			"eventName": "INSERT",
			"dynamodb": {
				"NewImage": {
					"PK": {"S": "baz"},
					"SK": {"S": "honk"}
				}
			}
		}
	]
}`

var testDynamoDBRecords = []events.DynamoDBEventRecord{
	events.DynamoDBEventRecord{
		EventID:   "ev1",
		EventName: "INSERT",
		Change: events.DynamoDBStreamRecord{
			NewImage: map[string]events.DynamoDBAttributeValue{
				"PK": events.NewStringAttribute("foo"),
				"SK": events.NewStringAttribute("bonk"),
			},
		},
	},
	events.DynamoDBEventRecord{
		EventID:   "ev2",
		EventName: "INSERT",
		Change: events.DynamoDBStreamRecord{
			NewImage: map[string]events.DynamoDBAttributeValue{
				"PK": events.NewStringAttribute("bar"),
				"SK": events.NewStringAttribute("quux"),
			},
		},
	},
	events.DynamoDBEventRecord{
		EventID:   "ev3",
		EventName: "INSERT",
		Change: events.DynamoDBStreamRecord{
			NewImage: map[string]events.DynamoDBAttributeValue{
				"PK": events.NewStringAttribute("baz"),
				"SK": events.NewStringAttribute("honk"),
			},
		},
	},
}

func (s *DynamoDBSuite) TestEventInit(t sweet.T) {
	handler := NewMockDynamoDBEventHandlerInitializer()
	outer := &dynamoDBEventHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: nacelle.NewServiceContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err).To(BeNil())
	Expect(handler.InitFunc).To(BeCalledOnceWith(config))
}

func (s *DynamoDBSuite) TestEventBadInjection(t sweet.T) {
	handler := &badInjectionDynamoDBEventHandler{}
	outer := &dynamoDBEventHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err.Error()).To(ContainSubstring("ServiceA"))
}

func (s *DynamoDBSuite) TestEventInitError(t sweet.T) {
	handler := NewMockDynamoDBEventHandlerInitializer()
	handler.InitFunc.SetDefaultReturn(fmt.Errorf("oops"))
	outer := &dynamoDBEventHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: nacelle.NewServiceContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err).To(MatchError("oops"))
}

func (s *DynamoDBSuite) TestRecordInit(t sweet.T) {
	handler := NewMockDynamoDBRecordHandlerInitializer()
	outer := &dynamoDBRecordHandler{
		handler:  handler,
		Services: nacelle.NewServiceContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err).To(BeNil())
	Expect(handler.InitFunc).To(BeCalledOnceWith(config))
}

func (s *DynamoDBSuite) TestRecordBadInjection(t sweet.T) {
	handler := &badInjectionDynamoDBRecordHandler{}
	outer := &dynamoDBRecordHandler{
		handler:  handler,
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err.Error()).To(ContainSubstring("ServiceA"))
}

func (s *DynamoDBSuite) TestRecordInitError(t sweet.T) {
	handler := NewMockDynamoDBRecordHandlerInitializer()
	handler.InitFunc.SetDefaultReturn(fmt.Errorf("oops"))
	outer := &dynamoDBRecordHandler{
		handler:  handler,
		Services: nacelle.NewServiceContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err).To(MatchError("oops"))
}

func (s *DynamoDBSuite) TestEventInvoke(t sweet.T) {
	handler := NewMockDynamoDBEventHandlerInitializer()
	outer := &dynamoDBEventHandler{
		handler: handler,
		Logger:  nacelle.NewNilLogger(),
	}

	response, err := outer.Invoke(context.Background(), []byte(testDynamoDBPayload))
	Expect(err).To(BeNil())
	Expect(response).To(BeNil())
	Expect(handler.HandleFunc).To(BeCalledOnceWith(BeAnything(), testDynamoDBRecords, BeAnything()))
}

func (s *DynamoDBSuite) TestEventInvokeError(t sweet.T) {
	handler := NewMockDynamoDBEventHandlerInitializer()
	outer := &dynamoDBEventHandler{
		handler: handler,
		Logger:  nacelle.NewNilLogger(),
	}

	handler.HandleFunc.SetDefaultReturn(fmt.Errorf("oops"))
	_, err := outer.Invoke(context.Background(), []byte(testDynamoDBPayload))
	Expect(err).To(MatchError("failed to process DynamoDB event (oops)"))
}

func (s *DynamoDBSuite) TestRecordHandle(t sweet.T) {
	handler := NewMockDynamoDBRecordHandlerInitializer()
	outer := &dynamoDBRecordHandler{handler: handler}

	err := outer.Handle(context.Background(), testDynamoDBRecords, nacelle.NewNilLogger())
	Expect(err).To(BeNil())

	for _, record := range testDynamoDBRecords {
		Expect(handler.HandleFunc).To(BeCalledOnceWith(BeAnything(), record, BeAnything()))
	}
}

func (s *DynamoDBSuite) TestRecordHandleError(t sweet.T) {
	handler := NewMockDynamoDBRecordHandlerInitializer()
	handler.HandleFunc.PushReturn(nil)
	handler.HandleFunc.PushReturn(fmt.Errorf("oops"))
	outer := &dynamoDBRecordHandler{handler: handler}

	err := outer.Handle(context.Background(), testDynamoDBRecords, nacelle.NewNilLogger())
	Expect(err).To(MatchError("failed to process DynamoDB record ev2 (oops)"))
	Expect(handler.HandleFunc).To(BeCalledN(2))
}

//
// Bad Injection

type badInjectionDynamoDBEventHandler struct {
	ServiceA *A `service:"A"`
}

func (i *badInjectionDynamoDBEventHandler) Handle(ctx context.Context, records []events.DynamoDBEventRecord, logger nacelle.Logger) error {
	return nil
}

type badInjectionDynamoDBRecordHandler struct {
	ServiceA *A `service:"A"`
}

func (i *badInjectionDynamoDBRecordHandler) Handle(ctx context.Context, record events.DynamoDBEventRecord, logger nacelle.Logger) error {
	return nil
}

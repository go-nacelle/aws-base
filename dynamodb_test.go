package lambdabase

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	mockassert "github.com/derision-test/go-mockgen/testutil/assert"
	"github.com/go-nacelle/nacelle"
	"github.com/stretchr/testify/require"
)

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
	{
		EventID:   "ev1",
		EventName: "INSERT",
		Change: events.DynamoDBStreamRecord{
			NewImage: map[string]events.DynamoDBAttributeValue{
				"PK": events.NewStringAttribute("foo"),
				"SK": events.NewStringAttribute("bonk"),
			},
		},
	},
	{
		EventID:   "ev2",
		EventName: "INSERT",
		Change: events.DynamoDBStreamRecord{
			NewImage: map[string]events.DynamoDBAttributeValue{
				"PK": events.NewStringAttribute("bar"),
				"SK": events.NewStringAttribute("quux"),
			},
		},
	},
	{
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

func TestDynamoDBEventInit(t *testing.T) {
	handler := NewMockDynamoDBEventHandlerInitializer()
	outer := &dynamoDBEventHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: nacelle.NewServiceContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	require.Nil(t, err)
	mockassert.CalledOnceWith(t, handler.InitFunc, mockassert.Values(config))
}

func TestDynamoDBEventBadInjection(t *testing.T) {
	handler := &badInjectionDynamoDBEventHandler{}
	outer := &dynamoDBEventHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "ServiceA")
}

func TestDynamoDBEventInitError(t *testing.T) {
	handler := NewMockDynamoDBEventHandlerInitializer()
	handler.InitFunc.SetDefaultReturn(fmt.Errorf("oops"))
	outer := &dynamoDBEventHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: nacelle.NewServiceContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	require.EqualError(t, err, "oops")
}

func TestDynamoDBRecordInit(t *testing.T) {
	handler := NewMockDynamoDBRecordHandlerInitializer()
	outer := &dynamoDBRecordHandler{
		handler:  handler,
		Services: nacelle.NewServiceContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	require.Nil(t, err)
	mockassert.CalledOnceWith(t, handler.InitFunc, mockassert.Values(config))
}

func TestDynamoDBRecordBadInjection(t *testing.T) {
	handler := &badInjectionDynamoDBRecordHandler{}
	outer := &dynamoDBRecordHandler{
		handler:  handler,
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "ServiceA")
}

func TestDynamoDBRecordInitError(t *testing.T) {
	handler := NewMockDynamoDBRecordHandlerInitializer()
	handler.InitFunc.SetDefaultReturn(fmt.Errorf("oops"))
	outer := &dynamoDBRecordHandler{
		handler:  handler,
		Services: nacelle.NewServiceContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	require.EqualError(t, err, "oops")
}

func TestDynamoDBEventInvoke(t *testing.T) {
	handler := NewMockDynamoDBEventHandlerInitializer()
	outer := &dynamoDBEventHandler{
		handler: handler,
		Logger:  nacelle.NewNilLogger(),
	}

	response, err := outer.Invoke(context.Background(), []byte(testDynamoDBPayload))
	require.Nil(t, err)
	require.Nil(t, response)
	mockassert.CalledOnceWith(t, handler.HandleFunc, mockassert.Values(mockassert.Skip, testDynamoDBRecords))
}

func TestDynamoDBEventInvokeError(t *testing.T) {
	handler := NewMockDynamoDBEventHandlerInitializer()
	outer := &dynamoDBEventHandler{
		handler: handler,
		Logger:  nacelle.NewNilLogger(),
	}

	handler.HandleFunc.SetDefaultReturn(fmt.Errorf("oops"))
	_, err := outer.Invoke(context.Background(), []byte(testDynamoDBPayload))
	require.EqualError(t, err, "failed to process DynamoDB event (oops)")
}

func TestDynamoDBRecordHandle(t *testing.T) {
	handler := NewMockDynamoDBRecordHandlerInitializer()
	outer := &dynamoDBRecordHandler{handler: handler}

	err := outer.Handle(context.Background(), testDynamoDBRecords, nacelle.NewNilLogger())
	require.Nil(t, err)

	for _, record := range testDynamoDBRecords {
		mockassert.CalledOnceWith(t, handler.HandleFunc, mockassert.Values(mockassert.Skip, record))
	}
}

func TestDynamoDBRecordHandleError(t *testing.T) {
	handler := NewMockDynamoDBRecordHandlerInitializer()
	handler.HandleFunc.PushReturn(nil)
	handler.HandleFunc.PushReturn(fmt.Errorf("oops"))
	outer := &dynamoDBRecordHandler{handler: handler}

	err := outer.Handle(context.Background(), testDynamoDBRecords, nacelle.NewNilLogger())
	require.EqualError(t, err, "failed to process DynamoDB record ev2 (oops)")
	mockassert.CalledN(t, handler.HandleFunc, 2)
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

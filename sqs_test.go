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

var testSQSPayload = `{
	"Records": [
		{"messageId": "m1", "body": "foo"},
		{"messageId": "m2", "body": "bar"},
		{"messageId": "m3", "body": "baz"}
	]
}`

var testSQSMessages = []events.SQSMessage{
	{MessageId: "m1", Body: "foo"},
	{MessageId: "m2", Body: "bar"},
	{MessageId: "m3", Body: "baz"},
}

func TestSQSEventInit(t *testing.T) {
	handler := NewMockSqsEventHandlerInitializer()
	outer := &sqsEventHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: nacelle.NewServiceContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	require.Nil(t, err)
	mockassert.CalledOnceWith(t, handler.InitFunc, mockassert.Values(config))
}

func TestSQSEventBadInjection(t *testing.T) {
	handler := &badInjectionSQSEventHandler{}
	outer := &sqsEventHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "ServiceA")
}

func TestSQSEventInitError(t *testing.T) {
	handler := NewMockSqsEventHandlerInitializer()
	handler.InitFunc.SetDefaultReturn(fmt.Errorf("oops"))
	outer := &sqsEventHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	require.EqualError(t, err, "oops")
}

func TestSQSMessageInit(t *testing.T) {
	handler := NewMockSqsMessageHandlerInitializer()
	outer := &sqsMessageHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: nacelle.NewServiceContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	require.Nil(t, err)
	mockassert.CalledOnceWith(t, handler.InitFunc, mockassert.Values(config))
}

func TestSQSMessageBadInjection(t *testing.T) {
	handler := &badInjectionSQSMessageHandler{}
	outer := &sqsMessageHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "ServiceA")
}

func TestSQSMessageInitError(t *testing.T) {
	handler := NewMockSqsMessageHandlerInitializer()
	handler.InitFunc.SetDefaultReturn(fmt.Errorf("oops"))
	outer := &sqsMessageHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	require.EqualError(t, err, "oops")
}

func TestSQSEventInvoke(t *testing.T) {
	handler := NewMockSqsEventHandlerInitializer()
	outer := &sqsEventHandler{
		handler: handler,
		Logger:  nacelle.NewNilLogger(),
	}

	response, err := outer.Invoke(context.Background(), []byte(testSQSPayload))
	require.Nil(t, err)
	require.Nil(t, response)
	mockassert.CalledOnceWith(t, handler.HandleFunc, mockassert.Values(mockassert.Skip, testSQSMessages))
}

func TestSQSEventInvokeError(t *testing.T) {
	handler := NewMockSqsEventHandlerInitializer()
	outer := &sqsEventHandler{
		handler: handler,
		Logger:  nacelle.NewNilLogger(),
	}

	handler.HandleFunc.SetDefaultReturn(fmt.Errorf("oops"))
	_, err := outer.Invoke(context.Background(), []byte(testSQSPayload))
	require.EqualError(t, err, "failed to process SQS event (oops)")
}

func TestSQSMessageHandle(t *testing.T) {
	handler := NewMockSqsMessageHandlerInitializer()
	outer := &sqsMessageHandler{handler: handler}

	err := outer.Handle(context.Background(), testSQSMessages, nacelle.NewNilLogger())
	require.Nil(t, err)

	for _, message := range testSQSMessages {
		mockassert.CalledOnceWith(t, handler.HandleFunc, mockassert.Values(mockassert.Skip, message))
	}
}

func TestSQSMessageHandleError(t *testing.T) {
	handler := NewMockSqsMessageHandlerInitializer()
	handler.HandleFunc.PushReturn(nil)
	handler.HandleFunc.PushReturn(fmt.Errorf("oops"))
	outer := &sqsMessageHandler{handler: handler}

	err := outer.Handle(context.Background(), testSQSMessages, nacelle.NewNilLogger())
	require.EqualError(t, err, "failed to process SQS message m2 (oops)")
	mockassert.CalledN(t, handler.HandleFunc, 2)
}

//
// Bad Injection

type badInjectionSQSEventHandler struct {
	ServiceA *A `service:"A"`
}

func (i *badInjectionSQSEventHandler) Handle(ctx context.Context, messages []events.SQSMessage, logger nacelle.Logger) error {
	return nil
}

type badInjectionSQSMessageHandler struct {
	ServiceA *A `service:"A"`
}

func (i *badInjectionSQSMessageHandler) Handle(ctx context.Context, message events.SQSMessage, logger nacelle.Logger) error {
	return nil
}

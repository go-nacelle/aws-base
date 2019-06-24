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

type SQSSuite struct{}

var testSQSPayload = `{
	"Records": [
		{"messageId": "m1", "body": "foo"},
		{"messageId": "m2", "body": "bar"},
		{"messageId": "m3", "body": "baz"}
	]
}`

var testSQSMessages = []events.SQSMessage{
	events.SQSMessage{MessageId: "m1", Body: "foo"},
	events.SQSMessage{MessageId: "m2", Body: "bar"},
	events.SQSMessage{MessageId: "m3", Body: "baz"},
}

func (s *SQSSuite) TestEventInit(t sweet.T) {
	handler := NewMockSqsEventHandlerInitializer()
	outer := &sqsEventHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: nacelle.NewServiceContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err).To(BeNil())
	Expect(handler.InitFunc).To(BeCalledOnceWith(config))
}

func (s *SQSSuite) TestEventBadInjection(t sweet.T) {
	handler := &badInjectionSQSEventHandler{}
	outer := &sqsEventHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err.Error()).To(ContainSubstring("ServiceA"))
}

func (s *SQSSuite) TestEventInitError(t sweet.T) {
	handler := NewMockSqsEventHandlerInitializer()
	handler.InitFunc.SetDefaultReturn(fmt.Errorf("oops"))
	outer := &sqsEventHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err).To(MatchError("oops"))
}

func (s *SQSSuite) TestMessageInit(t sweet.T) {
	handler := NewMockSqsMessageHandlerInitializer()
	outer := &sqsMessageHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: nacelle.NewServiceContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err).To(BeNil())
	Expect(handler.InitFunc).To(BeCalledOnceWith(config))
}

func (s *SQSSuite) TestMessageBadInjection(t sweet.T) {
	handler := &badInjectionSQSMessageHandler{}
	outer := &sqsMessageHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err.Error()).To(ContainSubstring("ServiceA"))
}

func (s *SQSSuite) TestMessageInitError(t sweet.T) {
	handler := NewMockSqsMessageHandlerInitializer()
	handler.InitFunc.SetDefaultReturn(fmt.Errorf("oops"))
	outer := &sqsMessageHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err).To(MatchError("oops"))
}

func (s *SQSSuite) TestEventInvoke(t sweet.T) {
	handler := NewMockSqsEventHandlerInitializer()
	outer := &sqsEventHandler{
		handler: handler,
		Logger:  nacelle.NewNilLogger(),
	}

	response, err := outer.Invoke(context.Background(), []byte(testSQSPayload))
	Expect(err).To(BeNil())
	Expect(response).To(BeNil())
	Expect(handler.HandleFunc).To(BeCalledOnceWith(BeAnything(), testSQSMessages, BeAnything()))
}

func (s *SQSSuite) TestEventInvokeError(t sweet.T) {
	handler := NewMockSqsEventHandlerInitializer()
	outer := &sqsEventHandler{
		handler: handler,
		Logger:  nacelle.NewNilLogger(),
	}

	handler.HandleFunc.SetDefaultReturn(fmt.Errorf("oops"))
	_, err := outer.Invoke(context.Background(), []byte(testSQSPayload))
	Expect(err).To(MatchError("failed to process SQS event (oops)"))
}

func (s *SQSSuite) TestMessageHandle(t sweet.T) {
	handler := NewMockSqsMessageHandlerInitializer()
	outer := &sqsMessageHandler{handler: handler}

	err := outer.Handle(context.Background(), testSQSMessages, nacelle.NewNilLogger())
	Expect(err).To(BeNil())

	for _, message := range testSQSMessages {
		Expect(handler.HandleFunc).To(BeCalledOnceWith(BeAnything(), message, BeAnything()))
	}
}

func (s *SQSSuite) TestMessageHandleError(t sweet.T) {
	handler := NewMockSqsMessageHandlerInitializer()
	handler.HandleFunc.PushReturn(nil)
	handler.HandleFunc.PushReturn(fmt.Errorf("oops"))
	outer := &sqsMessageHandler{handler: handler}

	err := outer.Handle(context.Background(), testSQSMessages, nacelle.NewNilLogger())
	Expect(err).To(MatchError("failed to process SQS message m2 (oops)"))
	Expect(handler.HandleFunc).To(BeCalledN(2))
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

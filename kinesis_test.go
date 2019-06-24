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

type KinesisSuite struct{}

var testKinesisPayload = `{
	"Records": [
		{
			"eventID": "ev1",
			"kinesis": {
				"PartitionKey": "foo",
				"Data": "WyJ4MSIsICJ5MSIsICJ6MSJdCg=="
			}
		},
		{
			"eventID": "ev2",
			"kinesis": {
				"PartitionKey": "bar",
				"Data": "WyJ4MiIsICJ5MiIsICJ6MiJdCg=="
			}
		},
		{
			"eventID": "ev3",
			"kinesis": {
				"PartitionKey": "baz",
				"Data": "WyJ4MyIsICJ5MyIsICJ6MyJdCg=="
			}
		}
	]
}`

var testKinesisRecords = []events.KinesisEventRecord{
	events.KinesisEventRecord{
		EventID: "ev1",
		Kinesis: events.KinesisRecord{
			PartitionKey: "foo",
			Data:         []byte{91, 34, 120, 49, 34, 44, 32, 34, 121, 49, 34, 44, 32, 34, 122, 49, 34, 93, 10},
		},
	},
	events.KinesisEventRecord{
		EventID: "ev2",
		Kinesis: events.KinesisRecord{
			PartitionKey: "bar",
			Data:         []byte{91, 34, 120, 50, 34, 44, 32, 34, 121, 50, 34, 44, 32, 34, 122, 50, 34, 93, 10},
		},
	},
	events.KinesisEventRecord{
		EventID: "ev3",
		Kinesis: events.KinesisRecord{
			PartitionKey: "baz",
			Data:         []byte{91, 34, 120, 51, 34, 44, 32, 34, 121, 51, 34, 44, 32, 34, 122, 51, 34, 93, 10},
		},
	},
}

func (s *KinesisSuite) TestEventInit(t sweet.T) {
	handler := NewMockKinesisEventHandlerInitializer()
	outer := &kinesisEventHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: nacelle.NewServiceContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err).To(BeNil())
	Expect(handler.InitFunc).To(BeCalledOnceWith(config))
}

func (s *KinesisSuite) TestEventBadInjection(t sweet.T) {
	handler := &badInjectionKinesisEventHandler{}
	outer := &kinesisEventHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err.Error()).To(ContainSubstring("ServiceA"))
}

func (s *KinesisSuite) TestEventInitError(t sweet.T) {
	handler := NewMockKinesisEventHandlerInitializer()
	handler.InitFunc.SetDefaultReturn(fmt.Errorf("oops"))
	outer := &kinesisEventHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err).To(MatchError("oops"))
}

func (s *KinesisSuite) TestRecordInit(t sweet.T) {
	handler := NewMockKinesisRecordHandlerInitializer()
	outer := &kinesisRecordHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: nacelle.NewServiceContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err).To(BeNil())
	Expect(handler.InitFunc).To(BeCalledOnceWith(config))
}

func (s *KinesisSuite) TestRecordBadInjection(t sweet.T) {
	handler := &badInjectionKinesisRecordHandler{}
	outer := &kinesisRecordHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err.Error()).To(ContainSubstring("ServiceA"))
}

func (s *KinesisSuite) TestRecordInitError(t sweet.T) {
	handler := NewMockKinesisRecordHandlerInitializer()
	handler.InitFunc.SetDefaultReturn(fmt.Errorf("oops"))
	outer := &kinesisRecordHandler{
		handler:  handler,
		Logger:   nacelle.NewNilLogger(),
		Services: makeBadContainer(),
	}

	config := nacelle.NewConfig(nacelle.NewTestEnvSourcer(nil))
	err := outer.Init(config)
	Expect(err).To(MatchError("oops"))
}

func (s *KinesisSuite) TestEventInvoke(t sweet.T) {
	handler := NewMockKinesisEventHandlerInitializer()
	outer := &kinesisEventHandler{
		handler: handler,
		Logger:  nacelle.NewNilLogger(),
	}

	response, err := outer.Invoke(context.Background(), []byte(testKinesisPayload))
	Expect(err).To(BeNil())
	Expect(response).To(BeNil())
	Expect(handler.HandleFunc).To(BeCalledOnceWith(BeAnything(), testKinesisRecords, BeAnything()))
}

func (s *KinesisSuite) TestEventInvokeError(t sweet.T) {
	handler := NewMockKinesisEventHandlerInitializer()
	outer := &kinesisEventHandler{
		handler: handler,
		Logger:  nacelle.NewNilLogger(),
	}

	handler.HandleFunc.SetDefaultReturn(fmt.Errorf("oops"))
	_, err := outer.Invoke(context.Background(), []byte(testKinesisPayload))
	Expect(err).To(MatchError("failed to process Kinesis event (oops)"))
}

func (s *KinesisSuite) TestRecordHandle(t sweet.T) {
	handler := NewMockKinesisRecordHandlerInitializer()
	outer := &kinesisRecordHandler{handler: handler}

	err := outer.Handle(context.Background(), testKinesisRecords, nacelle.NewNilLogger())
	Expect(err).To(BeNil())

	for _, record := range testKinesisRecords {
		Expect(handler.HandleFunc).To(BeCalledOnceWith(BeAnything(), record, BeAnything()))
	}
}

func (s *KinesisSuite) TestRecordHandleError(t sweet.T) {
	handler := NewMockKinesisRecordHandlerInitializer()
	handler.HandleFunc.PushReturn(nil)
	handler.HandleFunc.PushReturn(fmt.Errorf("oops"))
	outer := &kinesisRecordHandler{handler: handler}

	err := outer.Handle(context.Background(), testKinesisRecords, nacelle.NewNilLogger())
	Expect(err).To(MatchError("failed to process Kinesis record ev2 (oops)"))
	Expect(handler.HandleFunc).To(BeCalledN(2))
}

//
// Bad Injection

type badInjectionKinesisEventHandler struct {
	ServiceA *A `service:"A"`
}

func (i *badInjectionKinesisEventHandler) Handle(ctx context.Context, records []events.KinesisEventRecord, logger nacelle.Logger) error {
	return nil
}

type badInjectionKinesisRecordHandler struct {
	ServiceA *A `service:"A"`
}

func (i *badInjectionKinesisRecordHandler) Handle(ctx context.Context, record events.KinesisEventRecord, logger nacelle.Logger) error {
	return nil
}

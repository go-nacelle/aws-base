# Nacelle Base AWS Lambda Process [![GoDoc](https://godoc.org/github.com/go-nacelle/lambdabase?status.svg)](https://godoc.org/github.com/go-nacelle/lambdabase) [![CircleCI](https://circleci.com/gh/go-nacelle/lambdabase.svg?style=svg)](https://circleci.com/gh/go-nacelle/lambdabase) [![Coverage Status](https://coveralls.io/repos/github/go-nacelle/lambdabase/badge.svg?branch=master)](https://coveralls.io/github/go-nacelle/lambdabase?branch=master)

Abstract AWS Lambda server process for nacelle.

---

This library supplies an abstract AWS Lambda RPC server [process](https://nacelle.dev/docs/core/process) whose behavior can be be configured by implementing a `Handler` interface. This interface wraps the handler defined by [aws-lambda-go](https://github.com/aws/aws-lambda-go/blob/af0b813d5803d9754b920ed666b1cf8c16becfb3/lambda/handler.go#L14).

This library comes with an [example](https://github.com/go-nacelle/lambdabase/tree/master/example) project that logs values received from a Kinesis stream.

### Process

A Lambda server process is created by supplying a handler, described [below](https://nacelle.dev/docs/base-processes/lambdabase#handler), that controls its behavior.

```go
server := lambdabase.NewServer(NewHandler(), options...)
```

#### Event Sources

This library also supplies several additional abstract server processes that respond to specific Lambda [event sources](https://docs.aws.amazon.com/lambda/latest/dg/intro-invocation-modes.html). These servers require a more specific handler interface invoked with unmarshalled request data and additional log context.

<dl>
  <dt>NewDynamoDBEventServer</dt>
  <dd><a href="https://godoc.org/github.com/go-nacelle/lambdabase#NewDynamoDBEventServer">NewDynamoDBEventServer</a> invokes the backing handler with a list of DynamoDBEventRecords.</dd>

  <dt>NewDynamoDBRecordServer</dt>
  <dd><a href="https://godoc.org/github.com/go-nacelle/lambdabase#NewDynamoDBRecordServer">NewDynamoDBRecordServer</a> invokes the backing handler once for each DynamoDBEventRecord in the batch.</dd>

  <dt>NewKinesisEventServer</dt>
  <dd><a href="https://godoc.org/github.com/go-nacelle/lambdabase#NewKinesisEventServer">NewKinesisEventServer</a> invokes the backing handler with a list of KinesisEventRecords.</dd>

  <dt>NewKinesisRecordServer</dt>
  <dd><a href="https://godoc.org/github.com/go-nacelle/lambdabase#NewKinesisRecordServer">NewKinesisRecordServer</a> invokes the backing handler once for each KinesisEventRecord in the batch.</dd>

  <dt>NewSQSEventServer</dt>
  <dd><a href="https://godoc.org/github.com/go-nacelle/lambdabase#NewSQSEventServer">NewSQSEventServer</a> invokes the backing handler with a list of SQSMessages.</dd>

  <dt>NewSQSRecordServer</dt>
  <dd><a href="https://godoc.org/github.com/go-nacelle/lambdabase#NewSQSRecordServer">NewSQSRecordServer</a> invokes the backing handler once for each SQSMessage in the batch.</dd>
</dl>

### Handler

A handler is a struct with an `Init` and a `Handle` method. The initialization method, like the process that runs it, that takes a config object as a parameter. The handle method of the base server takes a context object and the request payload as parameters and returns the response payload and an error value. The handle method of an event-specific server takes a context object, the request payload, and a logger populated with request and event identifiers as parameters and returns an error value. Return an error from either method signals a fatal error to the process that runs it.

The following example implements a handler that transforms a request value to upper-case.

```go
type Handler struct {}

type ReqResp struct {
    Value string `json:"text"`
}

func (h *Handler) Init(config nacelle.Config) error {
    return nil
}

func (h *Handler) (ctx context.Context, payload []byte) ([]byte, error) {
    request := &ReqResp{}
    if err := json.Unmarshal(payload, &request); err != nil {
        return nil, err
    }

    return json.Marshal(&ReqResp{Value: strings.ToUpper(request.Value)})
}
```

### Configuration

The default process behavior can be configured by the following environment variables.

| Environment Variable | Required | Description |
| -------------------- | -------- | ----------- |
| _LAMBDA_SERVER_PORT  | yes      | The port on which to listen for RPC commands. |

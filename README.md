# Nacelle Base AWS Lambda Process [![GoDoc](https://godoc.org/github.com/go-nacelle/lambdabase?status.svg)](https://godoc.org/github.com/go-nacelle/lambdabase) [![CircleCI](https://circleci.com/gh/go-nacelle/lambdabase.svg?style=svg)](https://circleci.com/gh/go-nacelle/lambdabase) [![Coverage Status](https://coveralls.io/repos/github/go-nacelle/lambdabase/badge.svg?branch=master)](https://coveralls.io/github/go-nacelle/lambdabase?branch=master)

Abstract AWS Lambda server process for nacelle.

---

### Usage

The supplied server process is an abstract AWS Lambda RPC server whose behavior is determined by a supplied `Handler` interface which wraps the handler defined by [aws-lambda-go](https://github.com/aws/aws-lambda-go/blob/af0b813d5803d9754b920ed666b1cf8c16becfb3/lambda/handler.go#L14).

Supplied are six constructors supplied that create servers that specifically requests from an [event source mapping](https://docs.aws.amazon.com/lambda/latest/dg/intro-invocation-modes.html). These servers handle the request unmarshalling and add additional log context.

- **NewDynamoDBEventServer** invokes the backing handler with a list of DynamoDBEventRecords.
- **NewDynamoDBRecordServer** invokes the backing handler once for each DynamoDBEventRecord in the batch.
- **NewKinesisEventServer** invokes the backing handler with a list of KinesisEventRecords.
- **NewKinesisRecordServer** invokes the backing handler once for each KinesisEventRecord in the batch.
- **NewSQSEventServer** invokes the backing handler with a list of SQSMessages.
- **NewSQSRecordServer** invokes the backing handler once for each SQSMessage in the batch.

### Configuration

The default process behavior can be configured by the following environment variables.

| Environment Variable | Required | Description |
| -------------------- | -------- | ----------- |
| _LAMBDA_SERVER_PORT  | yes      | The port on which to listen for RPC commands. |

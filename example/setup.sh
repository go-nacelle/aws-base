#!/bin/bash -ex

LAMBDA_ENDPOINT=${LAMBDA_ENDPOINT:-http://localhost:4574}
KINESIS_ENDPOINT=${KINESIS_ENDPOINT:-http://localhost:4568}
STREAM_NAME=${STREAM_NAME:-lambdabase-test}
FUNCTION_NAME=${FUNCTION_NAME:-lambdabase-test}
STREAM_ARN="arn:aws:kinesis:us-east-1:000000000000:stream/${STREAM_NAME}"
STREAM_ARGS='--shard-count 1'
FUNCTION_ARGS='--handler example --role example --runtime go1.x --zip-file fileb://./example.zip'

GOOS=linux GOARCH=amd64 go build
zip example.zip example

aws kinesis create-stream --endpoint-url ${KINESIS_ENDPOINT} --stream-name ${STREAM_NAME} ${STREAM_ARGS}
aws lambda create-function --endpoint-url ${LAMBDA_ENDPOINT} --function-name ${FUNCTION_NAME} ${FUNCTION_ARGS}
aws lambda create-event-source-mapping --endpoint-url ${LAMBDA_ENDPOINT} --function-name ${FUNCTION_NAME} --event-source-arn ${STREAM_ARN}

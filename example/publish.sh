#!/bin/bash -ex

KINESIS_ENDPOINT=${KINESIS_ENDPOINT:-http://localhost:4568}
STREAM_NAME=${STREAM_NAME:-lambdabase-test}
MESSAGE=${MESSAGE:-`date "+%s"`}

aws kinesis put-record --endpoint-url ${KINESIS_ENDPOINT} --stream-name ${STREAM_NAME} --partition-key test --data "{\"data\": \"${MESSAGE}\"}"

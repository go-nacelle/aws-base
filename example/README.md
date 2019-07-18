# Lambda Base Example

A trivial example application to showcase the [lambdabase](https://nacelle.dev/docs/base-processes/lambdabase) library.

## Overview

This example application creates a Kinesis stream listener that prints its payloads.

## Building and Running

First, you must start [localstack](https://github.com/localstack/localstack). The provided docker-compose environment will start localstack with Kinesis and Lambda services exposed on ports 4568 and 4574, respectively.

```bash
docker-compose up localstack
```

After localstack has started, run the following bash script to compile the example binary, create a kinesis stream, create a lambda function with the binary, and create a kinesis event source.

```bash
./setup.sh
```

The lambda function will then print the payload of any kinesis message published (logs will be visible in the terminal running localstack). To publish a simple message, run the following bash script. If a message is not supplied, the current unix timestamp is sent as a payload.

```bash
MESSAGE='Hello, World!' ./publish.sh
```

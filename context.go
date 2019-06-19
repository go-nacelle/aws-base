package lambdabase

import (
	"context"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

func GetRequestID(ctx context.Context) string {
	if lc, ok := lambdacontext.FromContext(ctx); ok {
		return lc.AwsRequestID
	}

	return "<unknown request id>"
}

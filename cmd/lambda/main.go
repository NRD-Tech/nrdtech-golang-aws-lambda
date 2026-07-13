package main

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/lambda"

	"tmpl-go-lam/internal/loggingx"
)

func handler(ctx context.Context, event map[string]interface{}) (map[string]interface{}, error) {
	slog.Info("handling event")
	return map[string]interface{}{
		"statusCode": 200,
		"body":       "Hello World",
	}, nil
}

func main() {
	loggingx.ConfigureJSON()
	lambda.Start(handler)
}

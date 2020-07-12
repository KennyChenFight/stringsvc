package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"os"
)

func loggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			logger.Log("msg", "calling endpoint")
			defer logger.Log("msg", "called endpoint")
			return next(ctx, request)
		}
	}
}

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	svc := stringService{}

	var uppercase endpoint.Endpoint
	uppercase = makeUppercaseEndpoint(svc)
	uppercase = loggingMiddleware(log.With(logger, "method", "uppercase"))(uppercase)

	var count endpoint.Endpoint
	count = makeCountEndpoint(svc)
	count = loggingMiddleware(log.With(logger, "method", "count"))(count)

	uppercaseHandler := httptransport.NewServer(
		uppercase,
		decodeUppercaseRequest,
		encodeResponse)

	countHandler := httptransport.NewServer(
		count,
		decodeCountRequest,
		encodeResponse)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.ListenAndServe(":8080", nil)
}

package endpoint

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitprom "github.com/go-kit/kit/metrics/prometheus"
)

// MetricsMiddlewareHistorgram returns an endpoint.Endpoint that instruments endpoint latency using provided
// *kitprom.Histogram and LabelsFunc.
//
// Constructor panics if histogram or labels function are nil. endpoint.Middleware returned panics if next
// endpoint.Endpoint passed to it is nil.
func MetricsMiddlewareHistorgram(histogram *kitprom.Histogram, lf LabelsFunc) endpoint.Middleware {

	if histogram == nil {
		panic("metrics middleware counter: histogram required")
	}

	if lf == nil {
		panic("metrics middleware counter: labels func required")
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {

		if next == nil {
			panic("metrics middleware counter: next endpoint required")
		}

		return func(ctx context.Context, request any) (response any, err error) {
			defer func(begin time.Time) {
				histogram.With(lf(ctx, request, response, err)...).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)
		}
	}
}

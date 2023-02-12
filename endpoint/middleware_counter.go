package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	kitprom "github.com/go-kit/kit/metrics/prometheus"
)

// MetricsMiddlewareCounter returns an endpoint.Endpoint that instruments endpoint using provided *kitprom.Counter and
// LabelsFunc.
//
// Constructor panics if counter or labels function are nil. endpoint.Middleware returned panics if next
// endpoint.Endpoint passed to it is nil.
func MetricsMiddlewareCounter(counter *kitprom.Counter, lf LabelsFunc) endpoint.Middleware {

	if counter == nil {
		panic("metrics middleware counter: counter required")
	}

	if lf == nil {
		panic("metrics middleware counter: labels func required")
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {

		if next == nil {
			panic("metrics middleware counter: next endpoint required")
		}

		return func(ctx context.Context, request any) (response any, err error) {
			defer func() {
				counter.With(lf(ctx, request, response, err)...).Add(1)
			}()
			return next(ctx, request)
		}
	}
}

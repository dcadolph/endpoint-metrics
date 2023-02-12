package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	kitprom "github.com/go-kit/kit/metrics/prometheus"
)

type LabelsFunc func(ctx context.Context, req, resp any, err error) (labels []string)

func MetricsMiddlewareCounter(counter kitprom.Counter, lf LabelsFunc) endpoint.Middleware {

	if lf == nil {
		panic("metrics middleware counter: labels func required")
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {

		if next == nil {
			panic("metrics middleware counter: next endpoint required")
		}

		return func(ctx context.Context, request any) (response any, err error) {
			cf(ctx, )
		}
	}

}

package endpoint

import (
	"context"
	"io"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/go-kit/kit/endpoint"

	kitprom "github.com/go-kit/kit/metrics/prometheus"
	stdprom "github.com/prometheus/client_golang/prometheus"
)

func noOpLabelsFunc() LabelsFunc {
	return func(_ context.Context, _, _ any, _ error) []string {
		return nil
	}
}

func mockEndpoint(err error) endpoint.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		return nil, err
	}
}

// TestMetricsMiddlewareCounter_NilCounter tests that MetricsMiddlewareCounter constructor panics if counter passed to
// it is nil.
func TestMetricsMiddlewareCounter_NilCounter(t *testing.T) {

	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("run should have panicked but did not")
		}
	}()

	_ = MetricsMiddlewareCounter(nil, noOpLabelsFunc())
}

// TestMetricsMiddlewareCounter_NilLabelsFunc tests that MetricsMiddlewareCounter constructor panics if LabelsFunc
// passed to it is nil.
func TestMetricsMiddlewareCounter_NilLabelsFunc(t *testing.T) {

	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("run should have panicked but did not")
		}
	}()

	_ = MetricsMiddlewareCounter(&kitprom.Counter{}, nil)
}

// TestMetricsMiddlewareCounter_NilLabelsFunc tests that MetricsMiddlewareCounter panics if next endpoint.Endpoint is
// nil.
func TestMetricsMiddlewareCounter_NilNext(t *testing.T) {

	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("run should have panicked but did not")
		}
	}()

	_ = MetricsMiddlewareCounter(&kitprom.Counter{}, noOpLabelsFunc())(nil)
}

// TestMetricsMiddlewareCounter tests that MetricsMiddlewareCounter records the proper count and uses the right labels.
func TestMetricsMiddlewareCounter(t *testing.T) {

	vec := stdprom.NewCounterVec(
		stdprom.CounterOpts{
			Namespace: "space_case",
			Subsystem: "system",
			Name:      "request_total",
			Help:      "I need somebody...",
		},
		[]string{"error"},
	)

	stdprom.MustRegister(vec)

	counter := kitprom.NewCounter(vec)

	lf := func(_ context.Context, _, _ any, err error) []string {
		if err != nil {
			return []string{"error", "true"}
		}
		return []string{"error", "false"}
	}

	testTable := []struct {
		Endpoint   endpoint.Endpoint
		WantCount  int
		WantLabels map[string]string
	}{
		{
			Endpoint:   mockEndpoint(nil),
			WantCount:  1,
			WantLabels: map[string]string{"error": "false"},
		},
		{
			Endpoint:   mockEndpoint(io.EOF),
			WantCount:  1,
			WantLabels: map[string]string{"error": "true"},
		},
	}

	for testNum, test := range testTable {

		_, _ = MetricsMiddlewareCounter(counter, lf)(test.Endpoint)(context.Background(), "does-not-matter")

		if count := testutil.CollectAndCount(vec); count != test.WantCount {
			t.Fatalf("test %d: unexpected count: want: %d got: %d", testNum, test.WantCount, count)
		}

		if !vec.Delete(test.WantLabels) {
			t.Fatalf("test %d: expected labels not recorded", testNum)
		}
	}
}

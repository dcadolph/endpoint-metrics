package endpoint

import "context"

// LabelsFunc returns labels for instrumentation.
type LabelsFunc func(ctx context.Context, req, resp any, err error) (labels []string)

package app

import (
	"context"

	"github.com/bruli/pinger-exporter/internal/domain/metrics"
)

type ResourceMetricRepository interface {
	Create(ctx context.Context, r *metrics.Resource)
}

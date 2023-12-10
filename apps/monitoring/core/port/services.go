package port

import (
	"context"

	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/core/types/metrics"
)

type (
	MetricsService interface {
		GetMetrics() ([]metrics.Metric, error)
		ConfigureVisualizer(inst *containerstypes.Container) error
		InstallVisualizer(ctx context.Context, token string, visualizer string) error
		ConfigureCollector(inst *containerstypes.Container) error
		InstallCollector(ctx context.Context, token string, collector string) error
	}
)

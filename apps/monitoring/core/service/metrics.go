package service

import (
	"context"

	"github.com/google/uuid"
	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/core/port"
	"github.com/vertex-center/vertex/apps/monitoring/core/types/metrics"
)

type metricsService struct {
	uuid    uuid.UUID
	adapter port.MetricsAdapter
}

func NewMetricsService(metricsAdapter port.MetricsAdapter) port.MetricsService {
	return &metricsService{
		uuid:    uuid.New(),
		adapter: metricsAdapter,
	}
}

func (s *metricsService) GetMetrics() ([]metrics.Metric, error) {
	return s.adapter.GetMetrics()
}

func (s *metricsService) InstallCollector(ctx context.Context, token string, collector string) error {
	c := containersapi.NewContainersClient(token)

	serv, err := c.GetService(ctx, collector)
	if err != nil {
		return err
	}

	inst, err := c.InstallService(ctx, serv.ID)
	if err != nil {
		return err
	}

	err = s.ConfigureCollector(inst)
	if err != nil {
		return err
	}

	return c.PatchContainer(ctx, inst.UUID, containerstypes.ContainerSettings{
		Tags: []string{"Vertex Monitoring", "Vertex Monitoring - Prometheus Collector"},
	})
}

// ConfigureCollector will configure a container to monitor the metrics of Vertex.
func (s *metricsService) ConfigureCollector(inst *containerstypes.Container) error {
	return s.adapter.ConfigureContainer(inst.UUID)
}

func (s *metricsService) InstallVisualizer(ctx context.Context, token string, visualizer string) error {
	c := containersapi.NewContainersClient(token)

	serv, err := c.GetService(ctx, visualizer)
	if err != nil {
		return err
	}

	inst, err := c.InstallService(ctx, serv.ID)
	if err != nil {
		return err
	}

	err = s.ConfigureVisualizer(inst)
	if err != nil {
		return err
	}

	return c.PatchContainer(ctx, inst.UUID, containerstypes.ContainerSettings{
		Tags: []string{"Vertex Monitoring", "Vertex Monitoring - Grafana Visualizer"},
	})
}

func (s *metricsService) ConfigureVisualizer(inst *containerstypes.Container) error {
	// TODO: Implement
	return nil
}

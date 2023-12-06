package adapter

import (
	"context"

	"github.com/vertex-center/vertex/apps/admin/api"
	"github.com/vertex-center/vertex/apps/admin/core/port"
)

type HardwareApiAdapter struct{}

func NewHardwareApiAdapter() port.HardwareAdapter {
	return HardwareApiAdapter{}
}

func (HardwareApiAdapter) Reboot(ctx context.Context) error {
	return api.NewAdminKernelClient().Reboot(ctx)
}
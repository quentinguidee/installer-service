package containersapi

import (
	"context"

	"github.com/vertex-center/vertex/server/apps/containers/meta"
	"github.com/vertex-center/vertex/server/common/server"
	"github.com/vertex-center/vertex/server/config"
	"github.com/vertex-center/vertex/server/pkg/rest"
)

type Client struct {
	*rest.Client
}

func NewContainersClient(ctx context.Context) *Client {
	token := ctx.Value("token").(string)
	correlationID := ctx.Value(server.KeyCorrelationID).(string)
	return &Client{
		Client: rest.NewClient(config.Current.Addr(meta.Meta.ID), token, correlationID),
	}
}

type KernelClient struct {
	*rest.Client
}

func NewContainersKernelClient(ctx context.Context) *KernelClient {
	correlationID := ctx.Value(server.KeyCorrelationID).(string)
	return &KernelClient{
		Client: rest.NewClient(config.Current.KernelAddr(meta.Meta.ID), "", correlationID),
	}
}

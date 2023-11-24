package containersapi

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/pkg/router"
)

func (c *Client) GetContainer(ctx context.Context, uuid uuid.UUID) (*types.Container, *api.Error) {
	var inst types.Container
	var apiError api.Error

	err := c.Request().
		Pathf("./container/%s", uuid).
		ToJSON(&inst).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return &inst, api.HandleError(err, apiError)
}

func (c *Client) DeleteContainer(ctx context.Context, uuid uuid.UUID) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./container/%s", uuid).
		Delete().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func (c *Client) PatchContainer(ctx context.Context, uuid uuid.UUID, settings types.ContainerSettings) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./container/%s", uuid).
		Patch().
		BodyJSON(&settings).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func (c *Client) StartContainer(ctx context.Context, uuid uuid.UUID) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./container/%s/start", uuid).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func (c *Client) StopContainer(ctx context.Context, uuid uuid.UUID) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./container/%s/stop", uuid).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func (c *Client) PatchContainerEnvironment(ctx context.Context, uuid uuid.UUID, env map[string]string) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./container/%s/environment", uuid).
		Patch().
		BodyJSON(&env).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func (c *Client) GetDocker(ctx context.Context, uuid uuid.UUID) (map[string]any, *api.Error) {
	var info map[string]any
	var apiError api.Error
	err := c.Request().
		Pathf("./container/%s/docker", uuid).
		ToJSON(&info).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return info, api.HandleError(err, apiError)
}

func (c *Client) RecreateDocker(ctx context.Context, uuid uuid.UUID) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./container/%s/docker/recreate", uuid).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func (c *Client) GetContainerLogs(ctx context.Context, uuid uuid.UUID) (string, *api.Error) {
	var logs string
	var apiError api.Error
	err := c.Request().
		Pathf("./container/%s/logs", uuid).
		ToJSON(&logs).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return logs, api.HandleError(err, apiError)
}

func (c *Client) UpdateServiceContainer(ctx context.Context, uuid uuid.UUID) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./container/%s/update/service", uuid).
		Post().
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

func (c *Client) GetVersions(ctx context.Context, uuid uuid.UUID) ([]string, *api.Error) {
	var versions []string
	var apiError api.Error
	err := c.Request().
		Pathf("./container/%s/versions", uuid).
		ToJSON(&versions).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return versions, api.HandleError(err, apiError)
}

func (c *Client) WaitCondition(ctx context.Context, uuid uuid.UUID, condition container.WaitCondition) *api.Error {
	var apiError api.Error
	err := c.Request().
		Pathf("./container/%s/wait/%s", uuid, condition).
		ErrorJSON(&apiError).
		Fetch(ctx)
	return api.HandleError(err, apiError)
}

// Helpers

func GetContainerUUIDParam(c *router.Context) (uuid.UUID, *api.Error) {
	p := c.Param("container_uuid")
	if p == "" {
		return uuid.UUID{}, &api.Error{
			Code:    types.ErrCodeContainerUuidMissing,
			Message: "The request was missing the container UUID.",
		}
	}

	uid, err := uuid.Parse(p)
	if err != nil {
		return uuid.UUID{}, &api.Error{
			Code:    types.ErrCodeContainerUuidInvalid,
			Message: "The container UUID is invalid.",
		}
	}

	return uid, nil
}

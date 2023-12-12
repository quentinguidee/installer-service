package port

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/vertex-center/vertex/apps/containers/core/types"
)

type (
	ContainerAdapter interface {
		GetContainer(ctx context.Context, id types.ContainerID) (*types.Container, error)
		GetContainers(ctx context.Context) (types.Containers, error)
		CreateContainer(ctx context.Context, container types.Container) error
		DeleteContainer(ctx context.Context, id types.ContainerID) error
		GetTags(ctx context.Context) (types.Tags, error)
	}

	EnvAdapter interface {
		GetEnv(id types.ContainerID) (types.EnvVariables, error)
		SaveEnv(id types.ContainerID, env types.EnvVariables) error
	}

	LogsAdapter interface {
		Register(id types.ContainerID) error
		Unregister(id types.ContainerID) error
		UnregisterAll() error
		Push(id types.ContainerID, line types.LogLine)
		Pop(id types.ContainerID) (types.LogLine, error)
		LoadBuffer(id types.ContainerID) ([]types.LogLine, error) // LoadBuffer loads the latest logs kept in memory.
	}

	RunnerAdapter interface {
		DeleteContainer(ctx context.Context, c *types.Container) error
		DeleteMounts(ctx context.Context, c *types.Container) error
		Start(ctx context.Context, c *types.Container, setStatus func(status string)) (stdout io.ReadCloser, stderr io.ReadCloser, err error)
		Stop(ctx context.Context, c *types.Container) error
		Info(ctx context.Context, c types.Container) (map[string]any, error)
		WaitCondition(ctx context.Context, c *types.Container, cond types.WaitContainerCondition) error
		CheckForUpdates(ctx context.Context, c *types.Container) error
		HasUpdateAvailable(ctx context.Context, c types.Container) (bool, error)
		GetAllVersions(ctx context.Context, c types.Container) ([]string, error)
	}

	ServiceAdapter interface {
		Get(id string) (types.Service, error)
		GetRaw(id string) (interface{}, error)
		GetAll() []types.Service
		Reload() error
	}

	DockerAdapter interface {
		ListContainers() ([]types.DockerContainer, error)
		DeleteContainer(id string) error
		CreateContainer(options types.CreateContainerOptions) (types.CreateContainerResponse, error)
		StartContainer(id string) error
		StopContainer(id string) error
		InfoContainer(id string) (types.InfoContainerResponse, error)
		LogsStdoutContainer(id string) (io.ReadCloser, error)
		LogsStderrContainer(id string) (io.ReadCloser, error)
		WaitContainer(id string, cond types.WaitContainerCondition) error
		InfoImage(id string) (types.InfoImageResponse, error)
		PullImage(options types.PullImageOptions) (io.ReadCloser, error)
		BuildImage(options types.BuildImageOptions) (dockertypes.ImageBuildResponse, error)
	}
)

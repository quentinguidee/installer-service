package adapter

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/nat"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/juju/errors"
	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/pkg/vdocker"
	"github.com/vertex-center/vlog"
)

type runnerDockerAdapter struct{}

func NewRunnerDockerAdapter() port.RunnerAdapter {
	return runnerDockerAdapter{}
}

func (a runnerDockerAdapter) DeleteContainer(ctx context.Context, inst *types.Container) error {
	id, err := a.getContainerID(ctx, *inst)
	if err != nil {
		return err
	}

	cli := containersapi.NewContainersKernelClient(ctx)
	return cli.DeleteContainer(context.Background(), id)
}

func (a runnerDockerAdapter) DeleteMounts(ctx context.Context, inst *types.Container) error {
	id, err := a.getContainerID(ctx, *inst)
	if err != nil {
		return err
	}

	cli := containersapi.NewContainersKernelClient(ctx)
	return cli.DeleteMounts(context.Background(), id)
}

func (a runnerDockerAdapter) Start(ctx context.Context, c *types.Container, setStatus func(status string)) (io.ReadCloser, io.ReadCloser, error) {
	rErr, wErr := io.Pipe()
	rOut, wOut := io.Pipe()

	go func() {
		imageName := c.DockerImageVertexName()

		setStatus(types.ContainerStatusBuilding)

		log.Debug("building image", vlog.String("image", imageName))

		// Build
		stdout, err := a.buildImageFromName(ctx, c.GetImageNameWithTag())
		if err != nil {
			log.Error(err)
			setStatus(types.ContainerStatusError)
			return
		}

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer stdout.Close()

			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				var msg jsonmessage.JSONMessage
				err := json.Unmarshal(scanner.Bytes(), &msg)
				if err != nil {
					log.Error(err,
						vlog.String("text", scanner.Text()),
						vlog.String("uuid", c.ID.String()))
					continue
				}

				progress := types.DownloadProgress{
					ID:     msg.ID,
					Status: msg.Status,
				}

				if msg.Progress != nil {
					progress.Current = msg.Progress.Current
					progress.Total = msg.Progress.Total
				}

				progressJSON, err := json.Marshal(progress)
				if err != nil {
					log.Error(err,
						vlog.String("text", scanner.Text()),
						vlog.String("uuid", c.ID.String()))
					continue
				}

				_, err = fmt.Fprintf(wOut, "%s %s\n", "DOWNLOAD", progressJSON)
				if err != nil {
					log.Error(err,
						vlog.String("text", scanner.Text()),
						vlog.String("uuid", c.ID.String()))
					setStatus(types.ContainerStatusError)
					return
				}
			}
			if scanner.Err() != nil {
				log.Error(scanner.Err(),
					vlog.String("uuid", c.ID.String()))
				setStatus(types.ContainerStatusError)
				return
			}
		}()

		//wg.Add(1)
		//go func() {
		//	defer wg.Done()
		//	defer stderr.Close()
		//	_, err := io.Copy(wErr, stderr)
		//	if err != nil {
		//		log.Error(err)
		//		return
		//	}
		//}()

		log.Info("waiting for image to be built", vlog.String("uuid", c.ID.String()))

		wg.Wait()

		log.Info("image built", vlog.String("uuid", c.ID.String()))

		// Create
		id, err := a.getContainerID(ctx, *c)
		if errors.Is(err, errors.NotFound) {
			containerName := c.DockerContainerName()

			log.Info("container doesn't exists, create it.",
				vlog.String("container_name", containerName),
			)

			options := types.CreateContainerOptions{
				ContainerName: containerName,
				ExposedPorts:  nat.PortSet{},
				PortBindings:  nat.PortMap{},
				Binds:         []string{},
				Env:           []string{},
				CapAdd:        []string{},
			}

			// exposedPorts and portBindings
			var all []string
			for _, p := range c.Ports {
				for _, e := range c.Env {
					if e.Type == "port" && e.Name == p.Out {
						in := e.Value
						out := c.Env.Get(e.Name)
						all = append(all, out+":"+in)
						break
					}
				}
			}
			options.ExposedPorts, options.PortBindings, err = nat.ParsePortSpecs(all)
			if err != nil {
				return
			}

			// binds
			for _, v := range c.Volumes {
				in := v.In
				if !strings.HasPrefix(in, "/") {
					volumePath := a.getVolumePath(ctx, c.ID)
					in, err = filepath.Abs(path.Join(volumePath, in))
				}
				if err != nil {
					return
				}
				options.Binds = append(options.Binds, in+":"+v.Out)
			}

			// env
			for _, e := range c.Env {
				options.Env = append(options.Env, e.Name+"="+e.Value)
			}

			// capAdd
			for _, cap := range c.Caps {
				options.CapAdd = append(options.CapAdd, cap.Name)
			}

			// sysctls
			for _, sysctl := range c.Sysctls {
				options.Sysctls[sysctl.Name] = sysctl.Value
			}

			// cmd
			if c.Command != nil {
				options.Cmd = strings.Split(*c.Command, " ")
			}

			options.ImageName = c.GetImageNameWithTag()
			id, err = a.createContainer(ctx, options)
			if err != nil {
				return
			}
		} else if err != nil {
			return
		}

		// Start
		cli := containersapi.NewContainersKernelClient(ctx)
		err = cli.StartContainer(context.Background(), id)
		if err != nil {
			log.Error(err)
			setStatus(types.ContainerStatusError)
			return
		}
		setStatus(types.ContainerStatusRunning)

		var stderr io.ReadCloser
		stdout, stderr, err = a.readLogs(ctx, id)
		if err != nil {
			log.Error(err)
			return
		}

		go func() {
			defer stdout.Close()
			defer wOut.Close()

			_, err := io.Copy(wOut, stdout)
			if err != nil {
				log.Error(err)
				return
			}
		}()

		go func() {
			defer stderr.Close()
			defer wErr.Close()

			_, err := io.Copy(wOut, stdout)
			if err != nil {
				log.Error(err)
				return
			}
		}()

		err = a.WaitCondition(ctx, c, types.WaitContainerCondition(container.WaitConditionNotRunning))
		if err != nil {
			log.Error(err)
			setStatus(types.ContainerStatusError)
		} else {
			setStatus(types.ContainerStatusOff)
		}
	}()

	return rOut, rErr, nil
}

func (a runnerDockerAdapter) Stop(ctx context.Context, inst *types.Container) error {
	id, err := a.getContainerID(ctx, *inst)
	if err != nil {
		return err
	}

	cli := containersapi.NewContainersKernelClient(ctx)
	return cli.StopContainer(context.Background(), id)
}

func (a runnerDockerAdapter) Info(ctx context.Context, inst types.Container) (map[string]any, error) {
	id, err := a.getContainerID(ctx, inst)
	if err != nil {
		return nil, err
	}

	cli := containersapi.NewContainersKernelClient(ctx)
	info, err := cli.GetContainerInfo(context.Background(), id)
	if err != nil {
		return nil, err
	}

	imageInfo, err := cli.GetImageInfo(context.Background(), info.Image)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"container": info,
		"image":     imageInfo,
	}, nil
}

func (a runnerDockerAdapter) CheckForUpdates(ctx context.Context, inst *types.Container) error {
	imageName := inst.GetImageNameWithTag()

	res, err := a.pullImage(ctx, imageName)
	if err != nil {
		return err
	}
	defer res.Close()

	client := containersapi.NewContainersKernelClient(ctx)
	imageInfo, err := client.GetImageInfo(context.Background(), imageName)
	if err != nil {
		return err
	}

	latestImageID := imageInfo.ID

	currentImageID, err := a.getImageID(ctx, *inst)
	if err != nil {
		return err
	}

	if latestImageID == currentImageID {
		log.Info("already up-to-date",
			vlog.String("uuid", inst.ID.String()),
		)
		inst.Update = nil
	} else {
		log.Info("a new update is available",
			vlog.String("uuid", inst.ID.String()),
		)
		inst.Update = &types.ContainerUpdate{
			CurrentVersion: currentImageID,
			LatestVersion:  latestImageID,
		}
	}

	return nil
}

func (a runnerDockerAdapter) GetAllVersions(ctx context.Context, c types.Container) ([]string, error) {
	log.Debug("querying all versions of image", vlog.String("image", c.Image))
	return crane.ListTags(c.Image)
}

func (a runnerDockerAdapter) HasUpdateAvailable(ctx context.Context, inst types.Container) (bool, error) {
	//TODO implement me
	return false, nil
}

func (a runnerDockerAdapter) WaitCondition(ctx context.Context, inst *types.Container, cond types.WaitContainerCondition) error {
	id, err := a.getContainerID(ctx, *inst)
	if err != nil {
		return err
	}

	cli := containersapi.NewContainersKernelClient(ctx)
	return cli.WaitContainer(context.Background(), id, string(cond))
}

func (a runnerDockerAdapter) getContainer(ctx context.Context, inst types.Container) (types.DockerContainer, error) {
	cli := containersapi.NewContainersKernelClient(ctx)
	containers, err := cli.GetContainers(context.Background())
	if err != nil {
		return types.DockerContainer{}, err
	}

	var dockerContainer *types.DockerContainer
	for _, c := range containers {
		name := c.Names[0]
		if name == "/"+inst.DockerContainerName() {
			dockerContainer = &c
			break
		}
	}

	if dockerContainer == nil {
		return types.DockerContainer{}, errors.NotFoundf("docker container")
	}

	return *dockerContainer, nil
}

func (a runnerDockerAdapter) getContainerID(ctx context.Context, inst types.Container) (string, error) {
	c, err := a.getContainer(ctx, inst)
	if err != nil {
		return "", err
	}
	return c.ID, nil
}

func (a runnerDockerAdapter) getImageID(ctx context.Context, inst types.Container) (string, error) {
	c, err := a.getContainer(ctx, inst)
	if err != nil {
		return "", err
	}
	return c.ImageID, nil
}

func (a runnerDockerAdapter) pullImage(ctx context.Context, imageName string) (io.ReadCloser, error) {
	options := types.PullImageOptions{Image: imageName}

	cli := containersapi.NewContainersKernelClient(ctx)
	return cli.PullImage(context.Background(), options)
}

func (a runnerDockerAdapter) buildImageFromName(ctx context.Context, imageName string) (io.ReadCloser, error) {
	res, err := a.pullImage(ctx, imageName)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a runnerDockerAdapter) buildImageFromDockerfile(ctx context.Context, containerPath string, imageName string) (io.ReadCloser, error) {
	options := types.BuildImageOptions{
		Dir:        containerPath,
		Name:       imageName,
		Dockerfile: "Dockerfile",
	}

	cli := containersapi.NewContainersKernelClient(ctx)
	return cli.BuildImage(context.Background(), options)
}

func (a runnerDockerAdapter) createContainer(ctx context.Context, options types.CreateContainerOptions) (string, error) {
	cli := containersapi.NewContainersKernelClient(ctx)
	res, err := cli.CreateContainer(context.Background(), options)
	if err != nil {
		return "", err
	}

	for _, warn := range res.Warnings {
		log.Warn("warning while creating container",
			vlog.String("warning", warn),
		)
	}
	return res.ID, nil
}

func (a runnerDockerAdapter) readLogs(ctx context.Context, containerID string) (stdout io.ReadCloser, stderr io.ReadCloser, _ error) {
	cli := containersapi.NewContainersKernelClient(ctx)

	stdout, err := cli.GetContainerStdout(context.Background(), containerID)
	if err != nil {
		return nil, nil, err
	}

	stderr, err = cli.GetContainerStderr(context.Background(), containerID)
	if err != nil {
		stdout.Close()
		return nil, nil, err
	}

	return stdout, stderr, nil
}

func (a runnerDockerAdapter) getVolumePath(ctx context.Context, uuid types.ContainerID) string {
	appPath := a.getAppPath(ctx, "live_docker")
	return path.Join(appPath, "volumes", uuid.String())
}

func (a runnerDockerAdapter) getContainerPath(ctx context.Context, uuid types.ContainerID) string {
	appPath := a.getAppPath(ctx, "live")
	return path.Join(appPath, "containers", uuid.String())
}

func (a runnerDockerAdapter) getAppPath(ctx context.Context, base string) string {
	// If Vertex is running itself inside Docker, the containers are stored in the Vertex container volume.
	if vdocker.RunningInDocker() {
		cli := containersapi.NewContainersKernelClient(ctx)
		containers, err := cli.GetContainers(context.Background())
		if err != nil {
			log.Error(err)
		} else {
			for _, c := range containers {
				// find the docker container that has a volume /live, which is the Vertex container.
				for _, m := range c.Mounts {
					if m.Destination == "/"+base {
						base = m.Source
					}
				}
			}
		}
	}

	return path.Join(base, "apps", "containers")
}

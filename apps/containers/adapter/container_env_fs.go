package adapter

import (
	"bufio"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/storage"
)

const ContainerEnvPath = ".env"

type containerEnvFSAdapter struct {
	containersPath string
}

type ContainerEnvFSAdapterParams struct {
	containersPath string
}

func NewContainerEnvFSAdapter(params *ContainerEnvFSAdapterParams) port.ContainerEnvAdapter {
	if params == nil {
		params = &ContainerEnvFSAdapterParams{}
	}
	if params.containersPath == "" {
		params.containersPath = path.Join(storage.FSPath, "apps", "containers", "containers")
	}

	adapter := &containerEnvFSAdapter{
		containersPath: params.containersPath,
	}

	return adapter
}

func (a *containerEnvFSAdapter) Save(uuid types.ContainerID, env types.ContainerEnvVariables) error {
	envPath := path.Join(a.containersPath, uuid.String(), ContainerEnvPath)

	file, err := os.OpenFile(envPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	for key, value := range env {
		_, err := file.WriteString(strings.Join([]string{key, value}, "=") + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *containerEnvFSAdapter) Load(uuid types.ContainerID) (types.ContainerEnvVariables, error) {
	envPath := path.Join(a.containersPath, uuid.String(), ContainerEnvPath)

	file, err := os.Open(envPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	env := types.ContainerEnvVariables{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		if len(line) == 1 {
			env[line[0]] = ""
			continue
		}
		if len(line) == 2 {
			env[line[0]] = line[1]
			continue
		}
		return nil, errors.New("failed to read .env")
	}

	return env, nil
}

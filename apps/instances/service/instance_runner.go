package service

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/instances/adapter"
	"github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types/app"
	"github.com/vertex-center/vlog"
)

type InstanceRunnerService struct {
	ctx     *app.Context
	adapter types.InstanceRunnerAdapterPort
}

func NewInstanceRunnerService(ctx *app.Context, adapter types.InstanceRunnerAdapterPort) *InstanceRunnerService {
	return &InstanceRunnerService{
		ctx:     ctx,
		adapter: adapter,
	}
}

func (s *InstanceRunnerService) Install(uuid uuid.UUID, service types.Service) error {
	if service.Methods.Docker == nil {
		return ErrInstallMethodDoesNotExists
	}

	dir := path.Join(storage.Path, uuid.String())
	if service.Methods.Docker.Clone != nil {
		err := storage.CloneRepository(dir, service.Methods.Docker.Clone.Repository)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *InstanceRunnerService) Delete(inst *types.Instance) error {
	return s.adapter.Delete(inst)
}

// Start starts an instance by its UUID.
// If the instance does not exist, it returns ErrInstanceNotFound.
// If the instance is already running, it returns ErrInstanceAlreadyRunning.
func (s *InstanceRunnerService) Start(inst *types.Instance) error {
	if inst.IsBusy() {
		return nil
	}

	s.ctx.DispatchEvent(types.EventInstanceLog{
		InstanceUUID: inst.UUID,
		Kind:         types.LogKindOut,
		Message:      types.NewLogLineMessageString("Starting instance..."),
	})

	log.Info("starting instance",
		vlog.String("uuid", inst.UUID.String()),
	)

	if inst.IsRunning() {
		s.ctx.DispatchEvent(types.EventInstanceLog{
			InstanceUUID: inst.UUID,
			Kind:         types.LogKindVertexErr,
			Message:      types.NewLogLineMessageString(ErrInstanceAlreadyRunning.Error()),
		})
		return ErrInstanceAlreadyRunning
	}

	setStatus := func(status string) {
		s.setStatus(inst, status)
	}

	var runner types.InstanceRunnerAdapterPort
	if inst.IsDockerized() {
		runner = s.adapter
	} else {
		return fmt.Errorf("instance is not dockerized")
	}

	stdout, stderr, err := runner.Start(inst, setStatus)
	if err != nil {
		s.setStatus(inst, types.InstanceStatusError)
		return err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if scanner.Err() != nil {
				break
			}

			if strings.HasPrefix(scanner.Text(), "DOWNLOAD") {
				msg := strings.TrimPrefix(scanner.Text(), "DOWNLOAD")

				var downloadProgress types.DownloadProgress
				err := json.Unmarshal([]byte(msg), &downloadProgress)
				if err != nil {
					log.Error(err)
					continue
				}

				s.ctx.DispatchEvent(types.EventInstanceLog{
					InstanceUUID: inst.UUID,
					Kind:         types.LogKindDownload,
					Message:      types.NewLogLineMessageDownload(&downloadProgress),
				})
				continue
			}

			s.ctx.DispatchEvent(types.EventInstanceLog{
				InstanceUUID: inst.UUID,
				Kind:         types.LogKindOut,
				Message:      types.NewLogLineMessageString(scanner.Text()),
			})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			if scanner.Err() != nil {
				break
			}
			s.ctx.DispatchEvent(types.EventInstanceLog{
				InstanceUUID: inst.UUID,
				Kind:         types.LogKindErr,
				Message:      types.NewLogLineMessageString(scanner.Text()),
			})
		}
	}()

	// Wait for the instance until stopped
	wg.Wait()

	// Log stopped
	s.ctx.DispatchEvent(types.EventInstanceLog{
		InstanceUUID: inst.UUID,
		Kind:         types.LogKindVertexOut,
		Message:      types.NewLogLineMessageString("Stopping instance..."),
	})
	log.Info("stopping instance",
		vlog.String("uuid", inst.UUID.String()),
	)

	return nil
}

// Stop stops an instance by its UUID.
// If the instance does not exist, it returns ErrInstanceNotFound.
// If the instance is not running, it returns ErrInstanceNotRunning.
func (s *InstanceRunnerService) Stop(inst *types.Instance) error {
	if inst.IsBusy() {
		return nil
	}

	if !inst.IsRunning() {
		s.ctx.DispatchEvent(types.EventInstanceLog{
			InstanceUUID: inst.UUID,
			Kind:         types.LogKindVertexErr,
			Message:      types.NewLogLineMessageString(ErrInstanceNotRunning.Error()),
		})
		return ErrInstanceNotRunning
	}

	s.setStatus(inst, types.InstanceStatusStopping)

	var err error
	if inst.IsDockerized() {
		err = s.adapter.Stop(inst)
	} else {
		return fmt.Errorf("instance is not dockerized")
	}

	if err == nil {
		s.ctx.DispatchEvent(types.EventInstanceLog{
			InstanceUUID: inst.UUID,
			Kind:         types.LogKindVertexOut,
			Message:      types.NewLogLineMessageString("Instance stopped."),
		})

		log.Info("instance stopped",
			vlog.String("uuid", inst.UUID.String()),
		)

		s.setStatus(inst, types.InstanceStatusOff)
	} else {
		s.setStatus(inst, types.InstanceStatusRunning)
	}

	return err
}

func (s *InstanceRunnerService) GetDockerContainerInfo(inst types.Instance) (map[string]any, error) {
	if !inst.IsDockerized() {
		return nil, errors.New("instance is not using docker")
	}

	info, err := s.adapter.Info(inst)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (s *InstanceRunnerService) GetAllVersions(inst *types.Instance, useCache bool) ([]string, error) {
	if !useCache || len(inst.CacheVersions) == 0 {
		versions, err := s.adapter.GetAllVersions(*inst)
		if err != nil {
			return nil, err
		}
		inst.CacheVersions = versions
	}

	return inst.CacheVersions, nil
}

func (s *InstanceRunnerService) CheckForUpdates(inst *types.Instance) error {
	return s.adapter.CheckForUpdates(inst)
}

// RecreateContainer recreates a container by its UUID.
func (s *InstanceRunnerService) RecreateContainer(inst *types.Instance) error {
	if !inst.IsDockerized() {
		return nil
	}

	if inst.IsRunning() {
		err := s.adapter.Stop(inst)
		if err != nil {
			return err
		}
	}

	err := s.adapter.Delete(inst)
	if err != nil && !errors.Is(err, adapter.ErrContainerNotFound) {
		return err
	}

	go func() {
		err := s.Start(inst)
		if err != nil {
			log.Error(err)
			return
		}
	}()

	return nil
}

func (s *InstanceRunnerService) setStatus(inst *types.Instance, status string) {
	if inst.Status == status {
		return
	}

	var name string
	if inst.DisplayName == nil {
		name = inst.Service.Name
	} else {
		name = *inst.DisplayName
	}

	inst.Status = status
	s.ctx.DispatchEvent(types.EventInstancesChange{})
	s.ctx.DispatchEvent(types.EventInstanceStatusChange{
		InstanceUUID: inst.UUID,
		ServiceID:    inst.Service.ID,
		Instance:     *inst,
		Name:         name,
		Status:       status,
	})
}

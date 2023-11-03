package service

import (
	"context"
	"errors"
	"os"

	"github.com/vertex-center/vertex/apps/containers/core/types"
	vtypes "github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/event"

	"github.com/google/uuid"
	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	sqlapi "github.com/vertex-center/vertex/apps/sql/api"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

var (
	ErrPostgresDatabaseNotFound = errors.New("vertex postgres database not found")
)

type SetupService struct {
	uuid uuid.UUID
	ctx  *vtypes.VertexContext
}

func NewSetupService(ctx *vtypes.VertexContext) *SetupService {
	s := &SetupService{
		uuid: uuid.New(),
		ctx:  ctx,
	}
	s.ctx.AddListener(s)
	return s
}

func (s *SetupService) OnEvent(e event.Event) {
	switch e := e.(type) {
	case vtypes.EventAppReady:
		// TODO: The SQL app should also be ready!
		if e.AppID != "vx-containers" {
			return
		}
		go func() {
			err := s.setupDatabase()
			if err != nil {
				log.Error(err)
			}
		}()
	}
}

func (s *SetupService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *SetupService) setupDatabase() error {
	inst, err := s.getVertexDB()
	if err != nil && !errors.Is(err, ErrPostgresDatabaseNotFound) {
		return err
	}

	if inst == nil {
		err = s.installVertexDB()
		if err != nil {
			return err
		}

		inst, err = s.getVertexDB()
		if err != nil {
			return err
		}

		log.Info("vertex postgres database installed successfully",
			vlog.String("uuid", inst.UUID.String()))
	} else {
		log.Info("found vertex postgres container", vlog.String("uuid", inst.UUID.String()))
	}

	return s.startDatabase(inst)
}

func (s *SetupService) getVertexDB() (*types.Container, error) {
	insts, apiError := containersapi.GetContainers(context.Background())
	if apiError != nil {
		log.Error(apiError.RouterError())
		os.Exit(1)
	}

	for _, inst := range insts {
		isDatabase, isVertex, isPostgres := false, false, false
		if inst.Service.Features != nil && inst.Service.Features.Databases != nil {
			for _, db := range *inst.Service.Features.Databases {
				if db.Type == "postgres" {
					isPostgres = true
				}
			}
		}
		for _, tag := range inst.Tags {
			if tag == "Vertex SQL" {
				isDatabase = true
			}
			if tag == "Vertex Internal" {
				isVertex = true
			}
		}
		if isDatabase && isVertex && isPostgres {
			return inst, nil
		}
	}

	return nil, ErrPostgresDatabaseNotFound
}

func (s *SetupService) installVertexDB() error {
	log.Info("installing vertex postgres database")

	inst, apiError := sqlapi.InstallDBMS(context.Background(), "postgres")
	if apiError != nil {
		return apiError.RouterError()
	}

	inst.Tags = append(inst.Tags, "Vertex Internal")
	inst.DisplayName = "Vertex Database"

	apiError = containersapi.PatchContainer(context.Background(), inst.UUID, inst.ContainerSettings)
	if apiError != nil {
		return apiError.RouterError()
	}

	return nil
}

func (s *SetupService) startDatabase(inst *types.Container) error {
	eventsChan := make(chan interface{})
	defer close(eventsChan)

	abortChan := make(chan bool)
	defer close(abortChan)

	l := event.NewTempListener(func(e event.Event) {
		switch e := e.(type) {
		case types.EventContainerStatusChange:
			if e.ContainerUUID != inst.UUID {
				return
			}
			eventsChan <- e
		}
	})

	s.ctx.AddListener(l)
	defer s.ctx.RemoveListener(l)

	go func() {
		apiError := containersapi.StartContainer(context.Background(), inst.UUID)
		if apiError != nil {
			log.Error(apiError.RouterError())
		}
		abortChan <- true
	}()

	errFailedToStart := errors.New("failed to start vertex postgres database")

	for {
		select {
		case e := <-eventsChan:
			switch e := e.(type) {
			case types.EventContainerStatusChange:
				if e.Status == types.ContainerStatusRunning {
					return nil
				} else if e.Status == types.ContainerStatusOff || e.Status == types.ContainerStatusError {
					return errFailedToStart
				}
			}
		case <-abortChan:
			return errFailedToStart
		}
	}
}

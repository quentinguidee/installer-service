package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/juju/errors"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	ev "github.com/vertex-center/vertex/common/event"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/pkg/event"
)

func (s *containerService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *containerService) OnEvent(e event.Event) error {
	switch e := e.(type) {
	case ev.ServerSetupCompleted:
		go func() {
			err := s.StartAll(context.Background())
			if err != nil {
				log.Error(err)
			}
		}()
	case ev.ServerStop:
		return s.StopAll(context.Background())
	case types.EventContainerLog:
		return s.onLogReceived(e)
	}
	return nil
}

func (s *containerService) onLogReceived(e types.EventContainerLog) error {
	switch e.Kind {
	case types.LogKindDownload:
		var downloads *types.LogLineMessageDownloads
		download := e.Message.(*types.LogLineMessageDownload)

		line, err := s.logs.Pop(e.ContainerID)
		if err != nil && !errors.Is(err, types.ErrBufferEmpty) {
			return err
		}
		if line.Kind == types.LogKindDownloads {
			downloads = line.Message.(*types.LogLineMessageDownloads)
			downloads.Merge(download.DownloadProgress)
		} else {
			downloads = types.NewLogLineMessageDownloads(download.DownloadProgress)
		}
		s.logs.Push(e.ContainerID, types.LogLine{
			Kind:    types.LogKindDownloads,
			Message: downloads,
		})
	default:
		s.logs.Push(e.ContainerID, types.LogLine{
			Kind:    e.Kind,
			Message: e.Message,
		})
	}
	return nil
}
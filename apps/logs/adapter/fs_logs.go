package adapter

import (
	"sync"

	"github.com/vertex-center/vertex/apps/logs/core/port"
	"github.com/vertex-center/vlog"
)

type FSLogsAdapter struct {
	mu     sync.Mutex
	logger *vlog.Logger
}

func NewFSLogsAdapter() port.LogsAdapter {
	logger := vlog.New(
		vlog.WithOutputStd(),
		vlog.WithOutputFile(vlog.LogFormatText, "logs"),
	)
	a := &FSLogsAdapter{
		logger: logger,
	}
	var err error
	if err != nil {
		panic(err)
	}
	return a
}

func (a *FSLogsAdapter) Push(content string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.logger.Raw(content)
	return nil
}

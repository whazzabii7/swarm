package tasker

import (
	"github.com/whazzabii7/rpr"
	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/ui"
)

type TaskRequest rpr.RequestType
type Request = rpr.Request[TaskRequest]

const (
	LoadTask TaskRequest = iota + 300
)

type TaskManager struct {
	mfRequest   models.MFSubmit
	requestChan chan *Request
}

func NewTaskManager(requests models.MFSubmit) *TaskManager {
	return &TaskManager{
		mfRequest:   requests,
		requestChan: make(chan *Request, 100),
	}
}

func (t *TaskManager) Start() {
	for req := range t.requestChan {
		switch req.Type {
		default:
		}
	}
}

func (t *TaskManager) Stop(isStopped chan bool) {
	close(t.requestChan)
	ui.Log(ui.LevelInfo, "Tasker", "Stopped.")
	isStopped <- true
}

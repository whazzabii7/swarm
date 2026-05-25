package tasker

import (
	"github.com/whazzabii7/swarm/internal/rpr"
	"github.com/whazzabii7/swarm/internal/ui"
)

type TaskRequest rpr.RequestType
type Request = rpr.Request[TaskRequest]

const (
	LoadTask TaskRequest = iota + 300
)

type TaskManager struct {
	mfRequest   rpr.MFSubmit
	requestChan chan *Request
}

func NewTaskManager(requests rpr.MFSubmit) *TaskManager {
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

package tasker

import (
	"github.com/whazzabii7/swarm/internal/models" 
	"github.com/whazzabii7/swarm/internal/ui" 
)

type TaskRequest models.RequestType

const (
	LoadTask TaskRequest = iota + 300
)

type TaskManager struct {
	mfRequest chan models.Request[models.MFRequest]
	requestChan chan models.Request[TaskRequest]
}

func NewTaskManager(requests chan models.Request[models.MFRequest]) *TaskManager {
	return &TaskManager{
		mfRequest: requests,
		requestChan: make(chan models.Request[TaskRequest], 100),
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

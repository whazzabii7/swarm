package tasker

import (
	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/rpr"
	"github.com/whazzabii7/swarm/internal/ui"
)

type TaskRequest rpr.RequestType

const (
	LoadTask TaskRequest = iota + 300
)

type TaskManager struct {
	mfRequest   chan rpr.Request[rpr.MFRequest]
	requestChan chan rpr.Request[TaskRequest]
}

func NewTaskManager(requests chan rpr.Request[rpr.MFRequest]) *TaskManager {
	return &TaskManager{
		mfRequest:   requests,
		requestChan: make(chan rpr.Request[TaskRequest], 100),
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

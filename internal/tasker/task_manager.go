package tasker

import (
	"fmt"
	"github.com/whazzabii7/swarm/internal/models" 
)

type TaskRequest struct {}

type TaskManager struct {
	mfRequest chan models.MFRequest
	requestChan chan TaskRequest
}

func NewTaskManager(requests chan models.MFRequest) *TaskManager {
	return &TaskManager{
		mfRequest: requests,
		requestChan: make(chan TaskRequest),
	}
}

func (t *TaskManager) Start() {
	for req := range t.requestChan {
		switch req {
		default:
		}
	}
}

func (t *TaskManager) Stop() {
	close(t.requestChan)
	fmt.Println("[Tasker] Stopped.")
}

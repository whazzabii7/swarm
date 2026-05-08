package bot

import (
	"fmt"
	"github.com/whazzabii7/swarm/internal/models" 
)

type BotRequest struct {}

type BotManager struct {
	mfRequest chan models.MFRequest
	requestChan chan BotRequest
}

func NewBotManager(requests chan models.MFRequest) *BotManager {
	return &BotManager {
		mfRequest: requests,
		requestChan: make(chan BotRequest),
	}
}

func (b *BotManager) Start() {
	for req := range b.requestChan {
		switch req {
		default:
		}
	}
}

func (b *BotManager) Stop() {
	close(b.requestChan)
	fmt.Println("[BotManager] Stopped.")
}

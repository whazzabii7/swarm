package bot

import (
	"fmt"
	"github.com/whazzabii7/swarm/internal/models" 
)

type BotRequest struct {}

type ListenerMessage struct {}

type BotManager struct {
	mfRequest chan models.MFRequest
	requestChan chan BotRequest
	listenerChan chan ListenerMessage
	requestListener RequestListener
	botListener BotListener
}

func NewBotManager(requests chan models.MFRequest) *BotManager {
	rc := make(chan BotRequest)
	lc := make(chan ListenerMessage)
	return &BotManager {
		mfRequest: requests,
		requestChan: rc,
		listenerChan: lc,
		requestListener: *NewRequestListener(rc, lc),
		botListener: *NewBotListener(lc),
	}
}

func (b *BotManager) Start() {
	for req := range b.listenerChan {
		switch req {
		default:
		}
	}
}

func (b *BotManager) Stop() {
	b.requestListener.Stop()
	b.botListener.Stop()
	close(b.listenerChan)
	fmt.Println("[BotManager] Stopped.")
}

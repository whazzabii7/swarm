package bot

import (
	"context"
	"fmt"

	"github.com/whazzabii7/swarm/internal/models" 
)

type BotRequest models.RequestType

const (
	PingRequest BotRequest = iota + 200
)

type ListenerMessage struct {}

type BotManager struct {
	mfRequest chan models.Request[models.MFRequest] // <-chan, for sending requests to Mainframe
	requestChan chan models.Request[BotRequest]     // chan<-, for mainfrfame access to requestListener
	listenerChan chan ListenerMessage               // chan<-, for getting requests
	requestListener RequestListener                 // listens to requests from Mainframe
	botListener BotListener				            // listens to requests from Bots
}

func NewBotManager(requests chan models.Request[models.MFRequest]) *BotManager {
	rc := make(chan models.Request[BotRequest])
	lc := make(chan ListenerMessage)
	return &BotManager {
		mfRequest: requests,
		requestChan: rc,
		listenerChan: lc,
		requestListener: *NewRequestListener(rc, lc),
		botListener: *NewBotListener(lc),
	}
}

func (b *BotManager) Start(ctx context.Context, isStarted chan bool) {
	isSubStarted := make(chan bool)
	go b.requestListener.Start(ctx, isSubStarted)
	b.wait(isSubStarted)
	go b.botListener.Start(ctx, isSubStarted)
	b.wait(isSubStarted)
	isStarted<-true

	for req := range b.listenerChan {
		switch req {
		default:
		}
	}
}

func (b *BotManager) wait(cond chan bool) {
	if <-cond {
		return
	}
}

func (b *BotManager) Stop() {
	b.requestListener.Stop()
	b.botListener.Stop()
	close(b.listenerChan)
	fmt.Println("[BotManager] Stopped.")
}

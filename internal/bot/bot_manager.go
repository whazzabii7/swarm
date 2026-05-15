package bot

import (
	"context"
	"fmt"

	"github.com/whazzabii7/swarm/internal/models" 
)

type ListenerType int
const (
	ListenToMFRequest ListenerType = iota
	ListenToBots	
)

type BotRequest models.RequestType

const (
	BRPingRequest BotRequest = iota + 200
	BRStartBot
	BRStopBot
	BRSyncBlueprints
)

type ListenerMessage struct {
	source ListenerType
	requestType any
	payload any
	response chan any
}

type BotManager struct {
	mfRequest chan models.Request[models.MFRequest] // <-chan, for sending requests to Mainframe
	requestChan chan models.Request[BotRequest]     // chan<-, for mainfrfame access to requestListener
	listenerChan chan ListenerMessage               // chan<-, for getting requests
	requestListener RequestListener                 // listens to requests from Mainframe
	botListener BotListener				            // listens to requests from Bots
}

func NewBotManager(requests chan models.Request[models.MFRequest]) *BotManager {
	rc := make(chan models.Request[BotRequest], 100)
	lc := make(chan ListenerMessage, 100)
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
		switch req.source {
		case ListenToMFRequest:
			b.handleMFRequest(ctx, req)
		case ListenToBots:
			b.handleBotRequest(req)
		}
	}
}

func (b *BotManager) handleMFRequest(ctx context.Context, msg ListenerMessage) {
	switch msg.requestType {
	case BRStartBot:
		b.StartBot(ctx, msg.payload.(models.BotBlueprint))
	case BRStopBot:	
    case BRPingRequest:
	case BRSyncBlueprints:
	}
}

func (b *BotManager) handleBotRequest(msg ListenerMessage) {}

func (b *BotManager) wait(cond chan bool) {
	if <-cond {
		return
	}
}

func (b *BotManager) requestMainframe(request models.MFRequest, response chan models.Response) {
	b.mfRequest <- models.Request[models.MFRequest]{ Type: request, Payload: nil, Response: response }
}

func (b *BotManager) Submit(t BotRequest, data any, response chan models.Response ) {
	b.requestChan<-models.Request[BotRequest]{ Type: t, Payload: data, Response: response }
}

func (b *BotManager) Stop() {
	b.requestListener.Stop()
	b.botListener.Stop()
	close(b.listenerChan)
	fmt.Println("[BotManager] Stopped.")
}

package bot

import (
	"context"

	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/rpr"
	"github.com/whazzabii7/swarm/internal/ui"
)

type ListenerType int

const (
	ListenToMFRequest ListenerType = iota
	ListenToBots
)

type BotRequest rpr.RequestType

const (
	BRPingRequest BotRequest = iota + 200
	BRStartBot
	BRStopBot
	BRSyncBlueprints
)

type ListenerMessage struct {
	source      ListenerType
	requestType any
	payload     rpr.Payload
	response    chan rpr.Response
}

type BotManager struct {
	mfRequest       rpr.MFSubmit                 // <-chan, for sending requests to Mainframe
	requestChan     chan rpr.Request[BotRequest] // chan<-, for mainfrfame access to requestListener
	listenerChan    chan ListenerMessage            // chan<-, for getting requests
	requestListener RequestListener                 // listens to requests from Mainframe
	botListener     BotListener                     // listens to requests from Bots
}

func NewManager(requests rpr.MFSubmit) *BotManager {
	rc := make(chan rpr.Request[BotRequest], 100)
	lc := make(chan ListenerMessage, 100)
	return &BotManager{
		mfRequest:       requests,
		requestChan:     rc,
		listenerChan:    lc,
		requestListener: *NewRequestListener(rc, lc),
		botListener:     *NewBotListener(lc),
	}
}

func (b *BotManager) Start(ctx context.Context, isStarted chan bool) {
	isSubStarted := make(chan bool)
	go b.requestListener.Start(ctx, isSubStarted)
	b.wait(isSubStarted)
	go b.botListener.Start(ctx, isSubStarted)
	b.wait(isSubStarted)
	isStarted <- true

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
		if getBlueprint, ok := rpr.UnwrapPayload[func() models.BotBlueprint](msg.payload); ok {
			bp := getBlueprint()
			instance, err := b.startBot(ctx, bp)
			responseData := rpr.NewResponse[models.BotInstance](*instance, err)
			responseData.Submit(msg.response)
			rpr.NewResponse[models.BotInstance](*instance, err).Submit(msg.response)
		}
	case BRStopBot:
	case BRPingRequest:
	case BRSyncBlueprints:
		if getFuncArgs, ok := rpr.UnwrapPayload[func() (string, []models.BotBlueprint)](msg.payload); ok {
			blueprints, err := b.syncBlueprints(getFuncArgs())
			rpr.NewResponse[[]models.BotBlueprint](blueprints, err).Submit(msg.response)
		}
	}
}

func (b *BotManager) handleBotRequest(msg ListenerMessage) {}

func (b *BotManager) wait(cond chan bool) {
	if <-cond {
		return
	}
}

func (b *BotManager) Submit(t BotRequest, data any, response chan rpr.Response) {
	b.requestChan <- rpr.NewRequest[BotRequest](t, data, response)
}

func (b *BotManager) Stop(isStopped chan bool) {
	b.requestListener.Stop()
	b.botListener.Stop()
	close(b.listenerChan)
	ui.Log(ui.LevelInfo, "BotManager", "Stopped.")
	isStopped <- true
}

package bot

import (
	"context"

	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/rpr"
	"github.com/whazzabii7/swarm/internal/ui"
)

type BotRequest rpr.RequestType

const (
	BRPingRequest BotRequest = iota + 200
	BRStartBot
	BRStopBot
	BRSyncBlueprints
)

type BotManager struct {
	mfRequest    rpr.MFSubmit                  // <-chan, for sending requests to Mainframe
	requestChan  chan *rpr.Request[BotRequest] // chan<-, for mainfrfame access to requestListener
	listenerChan chan *rpr.Request[BotRequest] // chan<-, for getting requests
	botListener  *BotListener                  // listens to requests from Bots
}

func NewManager(requests rpr.MFSubmit) *BotManager {
	rc := make(chan *rpr.Request[BotRequest], 100)
	return &BotManager{
		mfRequest:    requests,
		requestChan:  rc,
		listenerChan: rc,
		botListener:  NewBotListener(rc),
	}
}

func (b *BotManager) Start(ctx context.Context, isStarted chan bool) {
	isSubStarted := make(chan bool)
	go b.botListener.Start(ctx, isSubStarted)
	b.wait(isSubStarted)
	isStarted <- true

	for req := range b.listenerChan {
		switch req.Type {
		case BRStartBot:
			if getBlueprint, ok := rpr.UnwrapPayload[func() models.BotBlueprint](req.Payload); ok {
				instance, err := b.startBot(ctx, getBlueprint())
				rpr.NewResponse[models.BotInstance](*instance, err).Submit(req.Response)
			}
		case BRStopBot:
		case BRPingRequest:
		case BRSyncBlueprints:
			if getFuncArgs, ok := rpr.UnwrapPayload[func() (string, []models.BotBlueprint)](req.Payload); ok {
				blueprints, err := b.syncBlueprints(getFuncArgs())
				rpr.NewResponse[[]models.BotBlueprint](blueprints, err).Submit(req.Response)
			}
		}
		req.Release()
	}
}

func (b *BotManager) wait(cond chan bool) {
	if <-cond {
		return
	}
}

func (b *BotManager) Submit(t BotRequest, data any, response chan *rpr.Response) {
	b.requestChan <- rpr.NewRequest[BotRequest](t, data, response)
}

func (b *BotManager) Stop(isStopped chan bool) {
	b.botListener.Stop()
	close(b.listenerChan)
	ui.Log(ui.LevelInfo, "BotManager", "Stopped.")
	isStopped <- true
}

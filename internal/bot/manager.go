package bot

import (
	"context"
	"fmt"

	"github.com/whazzabii7/rpr"
	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/ui"
)

type BotRequest rpr.RequestType
type Request = rpr.Request[BotRequest]

const (
	BRPingRequest BotRequest = iota + 200
	BRStartBot
	BRStopBot
	BRSyncBlueprints
)

type BotManager struct {
	mfRequest    models.MFSubmit // <-chan, for sending requests to Mainframe
	requestChan  chan *Request   // chan<-, for mainfrfame access to requestListener
	listenerChan chan *Request   // chan<-, for getting requests
	botListener  *BotListener    // listens to requests from Bots
}

func NewManager(requests models.MFSubmit) *BotManager {
	rc := rpr.MakeRequestChan[BotRequest](100)
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

	for {
		req, ok := <-b.listenerChan
		if !ok {
			break
		}
		if req == nil {
			continue
		}

		// Auslagerung der Logik schützt vor Deep-Nesting
		b.routeRequest(ctx, req)
	}
}

func (b *BotManager) routeRequest(ctx context.Context, req *Request) {
	defer req.Release()

	switch req.Type {
	case BRStartBot:
		b.handleStartBot(ctx, req)
	case BRStopBot:
		// TODO
	case BRPingRequest:
		// TODO
	case BRSyncBlueprints:
		b.handleSyncBlueprints(req)
	}
}

func (b *BotManager) handleStartBot(ctx context.Context, req *Request) {
	var blueprint models.BotBlueprint
	if ok := rpr.Assign(req.Payload.Get1(), &blueprint); !ok {
		err := fmt.Errorf("%w: Problem with Arguments", ErrFailUnpackArgs)
		rpr.NewResponseErr(err).Submit(req.Response)
		return
	}

	instance, err := b.startBot(ctx, blueprint)
	rpr.NewResponse(rpr.Pack(instance), err).Submit(req.Response)
}

func (b *BotManager) handleSyncBlueprints(req *Request) {
	var alias string
	var blueprints []models.BotBlueprint

	pay1, pay2 := req.Payload.Get2()
	ok1 := rpr.Assign(pay1, &alias)
	ok2 := rpr.Assign(pay2, &blueprints)
	if !ok1 || !ok2 {
		err := fmt.Errorf("%w: Problem with Arguments", ErrFailUnpackArgs)
		rpr.NewResponseErr(err).Submit(req.Response)
		return
	}

	blueprints, err := b.syncBlueprints(alias, blueprints)
	rpr.NewResponse(rpr.Pack(&blueprints), err).Submit(req.Response)
}

func (b *BotManager) wait(cond chan bool) {
	if <-cond {
		return
	}
}

func (b *BotManager) Submit(t BotRequest, data rpr.Payload, response chan *rpr.Response) {
	rpr.NewRequest[BotRequest](t, data, response).Submit(b.requestChan)
}

func (b *BotManager) Stop(isStopped chan bool) {
	b.botListener.Stop()
	close(b.listenerChan)
	ui.Log(ui.LevelInfo, "BotManager", "Stopped.")
	isStopped <- true
}

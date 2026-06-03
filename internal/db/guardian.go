package db

import (
	"context"

	"github.com/whazzabii7/rpr"
	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/ui"
)

type DBRequest rpr.RequestType
type Request = rpr.Request[DBRequest]

const (
	DBSaveBlueprint DBRequest = iota + 100
	DBGetBlueprint
	DBCheckBlueprints
	DBUpdateBotStatus
	DBGetActiveTasks
	DBRegisterInstance
)

type Guardian struct {
	mfRequest    models.MFSubmit // <-chan, for sending requests to Mainframe
	requestChan chan *Request
}

func NewGuardian(requests models.MFSubmit) *Guardian {
	return &Guardian{
		mfRequest:    requests,
		requestChan: rpr.MakeRequestChan[DBRequest](100),
	}
}

func (g *Guardian) Start(ctx context.Context, isStarted chan bool) {
	ui.Log(ui.LevelInfo, "Guardian", "Startup finished. Ready for requests...")
	isStarted <- true

	for req := range g.requestChan {
		g.routeRequest(ctx, req)
	}
}

func (g *Guardian) routeRequest(ctx context.Context, req *Request) {
	defer req.Release()
	ui.Logf(ui.LevelInfo, "DB-Guardian", "Recieved Request: %v", req.Type)

	switch req.Type {
	case DBSaveBlueprint:
		var blueprint models.BotBlueprint
		if ok := rpr.Assign(req.Payload.Get1(), &blueprint); ok {
			err := g.processSaveBlueprint(ctx, blueprint)
			rpr.NewResponseErr(err).Submit(req.Response)
		}

	case DBGetBlueprint:
		var alias string
		if ok := rpr.Assign(req.Payload.Get1(), &alias); ok {
			blueprint, err := g.handleGetBlueprint(ctx, alias)
			rpr.NewResponse(rpr.Pack(blueprint), err).Submit(req.Response)
		}

	case DBCheckBlueprints:
		var blueprints []models.BotBlueprint
		if ok := rpr.Assign(req.Payload.Get1(), &blueprints); ok {
			res, err := g.handleCheckBlueprints(ctx, blueprints)
			rpr.NewResponse(rpr.Pack(&res), err).Submit(req.Response)
		}

	case DBRegisterInstance:
		var instance models.BotInstance
		if ok := rpr.Assign(req.Payload.Get1(), &instance); ok {
			g.handleRegisterInstance(ctx, instance)
		}
	}
}

func (g *Guardian) Stop(isStopped chan bool) {
	close(g.requestChan)
	ui.Log(ui.LevelInfo, "Guardian", "Stopped.")
	isStopped <- true
}

func (g *Guardian) Submit(t DBRequest, data rpr.Payload, response chan *rpr.Response) {
	rpr.NewRequest[DBRequest](t, data, response).Submit(g.requestChan)
}

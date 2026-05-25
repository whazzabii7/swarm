package db

import (
	"context"

	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/rpr"
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
	requestChan chan *Request
}

func NewGuardian() *Guardian {
	return &Guardian{
		requestChan: make(chan *Request, 100),
	}
}

func (g *Guardian) Start(ctx context.Context, isStarted chan bool) {
	ui.Log(ui.LevelInfo, "Guardian", "Startup finished. Ready for requests...")
	isStarted <- true

	for req := range g.requestChan {
		ui.Logf(ui.LevelInfo, "DB-Guardian", "Recieved Request: %v", req.Type)
		switch req.Type {
		case DBSaveBlueprint:
			if getBlueprint, ok := rpr.UnwrapPayload[func() models.BotBlueprint](req.Payload); ok {
				err := g.processSaveBlueprint(ctx, getBlueprint())
				rpr.NewResponseErr(err).Submit(req.Response)
			}
		case DBGetBlueprint:
			if getBlueprintAlias, ok := rpr.UnwrapPayload[func() string](req.Payload); ok {
				blueprint, err := g.handleGetBlueprint(ctx, getBlueprintAlias())
				rpr.NewResponse[models.BotBlueprint](*blueprint, err).Submit(req.Response)
			}
		case DBCheckBlueprints:
			if getBlueprints, ok := rpr.UnwrapPayload[func() []models.BotBlueprint](req.Payload); ok {
				blueprints, err := g.handleCheckBlueprints(ctx, getBlueprints())
				rpr.NewResponse[[]models.BotBlueprint](blueprints, err).Submit(req.Response)
			}
		case DBRegisterInstance:
			if getInstance, ok := rpr.UnwrapPayload[func() models.BotInstance](req.Payload); ok {
				g.handleRegisterInstance(ctx, getInstance())
			}
		}
		req.Release()
	}
}

func (g *Guardian) Stop(isStopped chan bool) {
	close(g.requestChan)
	ui.Log(ui.LevelInfo, "Guardian", "Stopped.")
	isStopped <- true
}

func (g *Guardian) Submit(t DBRequest, data any, response chan *rpr.Response) {
	g.requestChan <- rpr.NewRequest[DBRequest](t, data, response)
}

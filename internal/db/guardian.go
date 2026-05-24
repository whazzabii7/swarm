package db

import (
	"context"

	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/ui"
)

type DBRequest models.RequestType

const (
	DBSaveBlueprint DBRequest = iota + 100
	DBGetBlueprint
	DBCheckBlueprints
	DBUpdateBotStatus
	DBGetActiveTasks
	DBRegisterInstance
)

type Guardian struct {
	requestChan chan models.Request[DBRequest]
}

func NewGuardian() *Guardian {
	return &Guardian{
		requestChan: make(chan models.Request[DBRequest], 100),
	}
}

func (g *Guardian) Start(ctx context.Context, isStarted chan bool) {
	ui.Log(ui.LevelInfo, "Guardian", "Startup finished. Ready for requests...")
	isStarted <- true

	for req := range g.requestChan {
		ui.Logf(ui.LevelInfo, "DB-Guardian", "Recieved Request: %v", req)
		switch req.Type {
		case DBSaveBlueprint:
			if getBlueprint, ok := models.UnwrapPayload[func() models.BotBlueprint](req.Payload); ok {
				err := g.processSaveBlueprint(ctx, getBlueprint())
				models.NewResponseErr(err).Submit(req.Response)
			}
		case DBGetBlueprint:
			if getBlueprintAlias, ok := models.UnwrapPayload[func() string](req.Payload); ok {
				blueprint, err := g.handleGetBlueprint(ctx, getBlueprintAlias())
				models.NewResponse[models.BotBlueprint](*blueprint, err).Submit(req.Response)
			}
		case DBCheckBlueprints:
			if getBlueprints, ok := models.UnwrapPayload[func() []models.BotBlueprint](req.Payload); ok {
				g.handleCheckBlueprints(ctx, getBlueprints())
			}
		case DBRegisterInstance:
			if getInstance, ok := models.UnwrapPayload[func() models.BotInstance](req.Payload); ok {
				g.handleRegisterInstance(ctx, getInstance())
			}
		default:
			ui.Log(ui.LevelInfo, "Guardian", "Unknown request type recieved")
		}
	}
}

func (g *Guardian) Stop(isStopped chan bool) {
	close(g.requestChan)
	ui.Log(ui.LevelInfo, "Guardian", "Stopped.")
	isStopped <- true
}

func (g *Guardian) Submit(t DBRequest, data any, response chan models.Response) {
	g.requestChan <- models.NewRequest[DBRequest](t, data, response)
}

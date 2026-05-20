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
	isStarted<-true

	for req := range g.requestChan {
		switch req.Type {
		case DBSaveBlueprint:
		    if bp, ok := req.Payload.(models.BotBlueprint); ok {
				g.processSaveBlueprint(ctx, bp)
			}
		case DBGetBlueprint:
			if bp, ok := req.Payload.(string); ok {
				g.handleGetBlueprint(ctx, bp, req.Response)
			}
		case DBRegisterInstance:
			bi, ok := req.Payload.(models.BotInstance)
			if ok {
				g.handleRegisterInstance(ctx, bi)
			}
		default:
			ui.Log(ui.LevelInfo, "Guardian", "Unknown request type recieved")
		}
	}
}

func (g *Guardian) Stop() {
	close(g.requestChan)
	ui.Log(ui.LevelInfo, "Guardian", "Stopped.")
}

func (g *Guardian) Submit(t DBRequest, data any, response chan models.Response) {
	g.requestChan <- models.Request[DBRequest]{Type: t, Payload: data, Response: response}
}

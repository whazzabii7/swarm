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
	isStarted<-true

	for req := range g.requestChan {
		switch req.Type {
		case DBSaveBlueprint:
		    if bp, ok := req.Payload.GetBluePrint(); ok {
				g.processSaveBlueprint(ctx, bp)
			}
		case DBGetBlueprint:
			if bp, ok := req.Payload.GetString(); ok {
				g.handleGetBlueprint(ctx, bp, req.Response)
			}
		case DBCheckBlueprints:
			if bps, ok := req.Payload.GetBluePrints(); ok {
				g.handleCheckBlueprints(ctx, bps)
			}
		case DBRegisterInstance:
			if bi, ok := req.Payload.GetInstance(); ok {
				g.handleRegisterInstance(ctx, bi)
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
	g.requestChan <- models.Request[DBRequest]{Type: t, Payload: models.Payload(data), Response: response}
}

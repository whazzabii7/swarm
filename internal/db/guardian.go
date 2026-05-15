package db

import (
	"context"
	"fmt"
	
	"github.com/whazzabii7/swarm/internal/models"
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
	fmt.Println("[DB-Guardian] Startup finished. Ready for requests...")
	isStarted<-true

	for req := range g.requestChan {
		switch req.Type {
		case DBSaveBlueprint:
			bp, ok := req.Payload.(models.BotBlueprint)
			if ok {
				g.processSaveBlueprint(ctx, bp)
			}
		case DBGetBlueprint:
		case DBRegisterInstance:
			bi, ok := req.Payload.(models.BotInstance)
			if ok {
				g.handleRegisterInstance(ctx, bi)
			}
		default:
			fmt.Println("[DB-Guardian] Unknown request type recieved")
		}
	}
}

func (g *Guardian) Stop() {
	close(g.requestChan)
	fmt.Println("[Guardian] Stopped.")
}

func (g *Guardian) Submit(t DBRequest, data any, response chan models.Response) {
	g.requestChan <- models.Request[DBRequest]{Type: t, Payload: data, Response: response}
}

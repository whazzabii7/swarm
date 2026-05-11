package db

import (
	"context"
	"fmt"
	
	"github.com/whazzabii7/swarm/internal/models"
)

type DBRequest models.RequestType

const (
	SaveBlueprint DBRequest = iota + 100
	UpdateBotStatus
	GetActiveTasks
	RegisterInstance
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
		case SaveBlueprint:
			bp, ok := req.Payload.(models.BotBlueprint)
			if ok {
				g.processSaveBlueprint(ctx, bp)
			}
		case RegisterInstance:
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

func (g *Guardian) Publish(t DBRequest, data any) {
	g.requestChan <- models.Request[DBRequest]{Type: t, Payload: data, Response: nil}
}

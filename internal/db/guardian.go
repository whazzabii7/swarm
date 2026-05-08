package db

import (
	"fmt"
	
	"github.com/whazzabii7/swarm/internal/models"
)

type RequestType int

const (
	SaveBlueprint RequestType = iota
	UpdateBotStatus
	GetActiveTasks
	RegisterInstance
)

type DBRequest struct {
	Type RequestType
	Data any
	Response chan string 
}

type Guardian struct {
	requestChan chan DBRequest
}

func NewGuardian() *Guardian {
	return &Guardian{
		requestChan: make(chan DBRequest, 100),
	}
}

func (g *Guardian) Start() {
	fmt.Println("[DB-Guardian] Startup finished. Ready for requests...")

	for req := range g.requestChan {
		switch req.Type {
		case SaveBlueprint:
			bp, ok := req.Data.(models.BotBlueprint)
			if ok {
				g.processSaveBlueprint(bp)
			}
		case RegisterInstance:
			bi, ok := req.Data.(models.BotInstance)
			if ok {
				g.handleRegisterInstance(bi)
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

func (g *Guardian) Publish(t RequestType, data any) {
	g.requestChan <- DBRequest{Type: t, Data: data, Response: nil}
}

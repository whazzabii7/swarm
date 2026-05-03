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
)

type DBRequest struct {
	Type RequestType
	Data interface{}
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

}

func (g *Guardian) processSaveBlueprint(bp models.BotBlueprint) {

}

func (g *Guardian) Puplish(t RequestType, data interface{}) {

}

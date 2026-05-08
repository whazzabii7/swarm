package models

import (
	"time"
)

type BotBlueprint struct {
	Alias       string    `json:"alias"`
	Path        string    `json:"path"`
	Type        string    `json:"type"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	LastScan    time.Time `json:"last_scan"`
}

type BotInstance struct {
	ID          int       `json:"id"`
	Alias       string    `json:"alias"` 
	PID         int       `json:"pid"`
	Status      string    `json:"status"`    
	LastSeen    time.Time `json:"last_seen"`
}

type Task struct {
	ID          string    `json:"id"` 
	TargetBot   string    `json:"target_bot"` 
	Payload     string    `json:"payload"`
	Status      string    `json:"status"` 
	CreatedAt   time.Time `json:"created_at"`
}

// Mainframe Requests
type MFRequest struct {
	RequestType int      `json:"request_type"`
	Payload     any      `json:"payload"`
	ResultChan  chan any `json:"requsted_data"`
}

package models

import (
	"time"
)

type BotBlueprint struct {
	Alias       string       `json:"alias"`
	Path        string       `json:"path"`
	Type        string       `json:"type"`
	Version     string       `json:"version"`
	Description string       `json:"description"`
	LastScan    time.Time    `json:"last_scan"`
}

type BotInstance struct {
	ID          int          `json:"id"`
	Alias       string       `json:"alias"` 
	PID         int          `json:"pid"`
	Status      string       `json:"status"`    
	LastSeen    time.Time    `json:"last_seen"`
}

type Task struct {
	ID          string        `json:"id"` 
	TargetBot   string        `json:"target_bot"` 
	Payload     string        `json:"payload"`
	Status      string        `json:"status"` 
	CreatedAt   time.Time     `json:"created_at"`
}

type Response struct {
	Payload any
	Err 	error
}

type RequestType int
type RequestConstraint interface { ~int }

type MFRequest RequestType

const (
	MFDataRequest MFRequest = iota
)

type Request[T RequestConstraint] struct {
	Type        T             `json:"type"`
	Payload     any           `json:"payload"`
	Response    chan Response `json:"response"`
}


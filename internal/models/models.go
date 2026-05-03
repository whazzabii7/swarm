package models

import (
	"time"
)

// BotBlueprint repräsentiert die installierte Binary (der Bauplan)
type BotBlueprint struct {
	Alias    string    `json:"alias"`
	Path     string    `json:"path"`
	Type     string    `json:"type"`
	LastScan time.Time `json:"last_scan"`
}

// BotInstance repräsentiert einen laufenden Prozess basierend auf einem Blueprint
type BotInstance struct {
	ID       int       `json:"id"`
	Alias    string    `json:"alias"` // Dies ist logisch dein Foreign Key zum Blueprint
	PID      int       `json:"pid"`
	Status   string    `json:"status"`    // starting, active, zombie, dead
	LastSeen time.Time `json:"last_seen"`
}

// Task repräsentiert einen konkreten Auftrag für einen Bot
type Task struct {
	ID        string    `json:"id"` // Tipp: Nutze hier String, falls du UUIDs verwenden willst
	TargetBot string    `json:"target_bot"` // Logischer Foreign Key zu BotBlueprint.Alias
	Payload   string    `json:"payload"`
	Status    string    `json:"status"` // pending, processing, done, failed
	CreatedAt time.Time `json:"created_at"`
}

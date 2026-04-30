package models

import (
	"time"
)

// Bot repräsentiert eine installierte Einheit
type Bot struct {
	ID           int       `json:"id"`
	Alias        string    `json:"alias"`
	Path         string    `json:"path"`
	Type         string    `json:"type"`
	Status       string    `json:"status"` // active, offline, corrupted
	Meta         string    `json:"meta"`   // JSON-Header mit Capabilities
	LastCheck    time.Time `json:"last_check"`
}

// Task repräsentiert eine Mission in der Queue
type Task struct {
	ID           string    `json:"id"`
	BotAlias     string    `json:"bot_alias"`
	Payload      string    `json:"payload"`      // Parameter für den Bot
	Status       string    `json:"status"`       // pending, running, done, failed
	Priority     int       `json:"priority"`
	Dependency   string    `json:"dependency"`   // ID eines anderen Tasks
	RetryCount   int       `json:"retry_count"`
	MaxRetries   int       `json:"max_retries"`
	Timeout      int       `json:"timeout"`      // Sekunden bis Kill
	Result       string    `json:"result"`
	CreatedAt    time.Time `json:"created_at"`
}

// Event für das System-Logging
type Event struct {
	ID        int       `json:"id"`
	TaskID    string    `json:"task_id"`
	Origin    string    `json:"origin"` // "SYSTEM" oder Bot-Name
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

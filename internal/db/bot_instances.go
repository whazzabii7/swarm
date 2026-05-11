package db

import (
	"context"
	"log"
	"time"

	"github.com/whazzabii7/swarm/internal/models"
)

func (g *Guardian) handleRegisterInstance(ctx context.Context, inst models.BotInstance) {
	query := `
	INSERT INTO bot_instances (alias, pid, status, last_seen)
	VALUES (?, ?, ?, ?);`

	_, err := DB.ExecContext(ctx, query, inst.Alias, inst.PID, inst.Status, inst.LastSeen.Format(time.RFC3339))
	if err != nil {
		log.Printf("[-] Guardian Error: Could not register instance for %s: %v", inst.Alias, err)
		return
	}
	log.Printf("[+] Guardian: Instance of '%s' registered with PID %d", inst.Alias, inst.PID)
}

package db

import (
	"context"
	"fmt"
	"time"

	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/ui"
)

func (g *Guardian) handleRegisterInstance(ctx context.Context, inst models.BotInstance) error {
	query := `
	INSERT INTO bot_instances (alias, pid, status, last_seen)
	VALUES (?, ?, ?, ?);`

	_, err := DB.ExecContext(ctx, query, inst.Alias, inst.PID, inst.Status, inst.LastSeen.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("%w %s: %v", ErrRegisterFailed, inst.Alias, err)
	}
	ui.Logf(ui.LevelInfo, "DB-Guardian", "[+] Instance of '%s' registered with PID %d", inst.Alias, inst.PID)
	return nil
}

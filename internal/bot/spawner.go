package bot

import(
	"fmt"
	"context"
	"os/exec"
	"time"

	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/ui"
)

func (m *BotManager) startBot(ctx context.Context, bp models.BotBlueprint) (*models.BotInstance, error) {
	cmd := exec.CommandContext(ctx, bp.Path)

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start bot %s: %v", bp.Alias, err)
	}

	instance := models.BotInstance{
		Alias:  bp.Alias,
		PID:    cmd.Process.Pid,
		Status: "active",
		LastSeen: time.Now().UTC().Truncate(time.Second),
	}

	ui.Logf(ui.LevelInfo, "BotManager", "Bot %s started with PID %d\n", instance.Alias, instance.PID)

	go func(){
		cmd.Wait()
		ui.Logf(ui.LevelInfo, "BotManager", "Bot %s (PID %d) stopped.\n", instance.Alias, instance.PID)
	}()

	return &instance, nil
}

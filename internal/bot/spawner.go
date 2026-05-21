package bot

import(
	"fmt"
	"bufio"
	"context"
	"os/exec"
	"time"

	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/ui"
)

func (m *BotManager) startBot(ctx context.Context, bp models.BotBlueprint) (*models.BotInstance, error) {
	cmd := exec.CommandContext(ctx, bp.Path)

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %v", err)
	}
	_ = stdinPipe

	// --------------------------------------------

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start bot %s: %v", bp.Alias, err)
	}

	instance := models.BotInstance{
		Alias:    bp.Alias,
		PID:      cmd.Process.Pid,
		Status:   "active",
		LastSeen: time.Now().UTC().Truncate(time.Second),
	}

	ui.Logf(ui.LevelInfo, "BotManager", "Bot %s started with PID %d\n", instance.Alias, instance.PID)

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			ui.Logf(ui.LevelInfo, fmt.Sprintf("BOT:%s:LOG", bp.Alias), "%s", scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			text := scanner.Text()
			ui.Logf(ui.LevelInfo, fmt.Sprintf("BOT:%s:DATA", bp.Alias), "%s", text)
			
			// HIER kommt später die Brücke zum BotListener hin!
		}
	}()
	// -----------------------------------------------------

	go func(){
		cmd.Wait()
		ui.Logf(ui.LevelInfo, "BotManager", "Bot %s (PID %d) stopped.\n", instance.Alias, instance.PID)
	}()

	return &instance, nil
}

package bot

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/whazzabii7/swarm/internal/models"
)

type BotManager struct {}

func NewBotManager() *BotManager {
	return &BotManager {}
}

func (m *BotManager) getBotHeader(path string) ( *models.BotBlueprint, error ) { 
	cmd := exec.Command(path, "--swarm-info")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var bp models.BotBlueprint
	if err := json.Unmarshal(output, &bp); err != nil {
		return nil, err
	}

	bp.Path = path
	bp.LastScan = time.Now().UTC().Truncate(time.Second)
	return &bp, nil 
}

func (m *BotManager) SyncBlueprints() ( []models.BotBlueprint, error ) {
	var foundBlueprints []models.BotBlueprint
	files, err := os.ReadDir("./bots")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// ignoring folder and non binaries
		if file.IsDir() { continue }

		info, err := file.Info()
		if err != nil { continue }

		if !( strings.HasSuffix(file.Name(), ".exe") || info.Mode().Perm()&0111 != 0 ) {
			continue
		}
		path := filepath.Join("./bots", file.Name())

		blueprint, err := m.getBotHeader(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Header of %s couldn't be read: %v\n", file.Name(), err)
			continue
		}
		foundBlueprints = append(foundBlueprints, *blueprint)
	}
	return foundBlueprints, nil
}

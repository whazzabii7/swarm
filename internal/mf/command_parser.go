package mf

import (
	"os"
	"bufio"
	"strings"

	// "github.com/whazzabii7/swarm/internal/models" 
	"github.com/whazzabii7/swarm/internal/ui" 
)

type CommandParser struct {
	// requestChan	chan models.MFRequest
	CommandChan chan Command `json:"command_chan"`
}

func NewComandParser() *CommandParser {
	return &CommandParser{
		// requestChan: requests,
		CommandChan: make(chan Command, 100),
	}
}

func (c *CommandParser) RunShell() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() { break }
		input := c.parse(scanner.Text())
		switch input.Type {
		case CmdQuit:
			c.CommandChan <- *input
			return
		default:
			c.CommandChan <- *input
		}
	}
}

func (c * CommandParser) Stop(isStopped chan bool) {
	close(c.CommandChan)
	ui.Log(ui.LevelInfo, "Commander", "Stopped.")
	isStopped <- true
}

func (c *CommandParser) parse(cmdStr string) *Command {
	cmd := strings.Split(cmdStr, " ")
	return NewCommand(collectArgs(cmd))
}

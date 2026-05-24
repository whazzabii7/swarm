package command

import (
	"bufio"
	"os"
	"strings"

	// "github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/ui"
)

type Parser struct {
	// requestChan	chan models.MFRequest
	CommandChan chan Command `json:"command_chan"`
}

func NewParser() *Parser {
	return &Parser{
		// requestChan: requests,
		CommandChan: make(chan Command, 100),
	}
}

func (p *Parser) RunShell() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			break
		}
		input := p.parse(scanner.Text())
		switch input.Type {
		case Quit:
			p.CommandChan <- *input
			return
		default:
			p.CommandChan <- *input
		}
	}
}

func (p *Parser) Stop(isStopped chan bool) {
	close(p.CommandChan)
	ui.Log(ui.LevelInfo, "Commander", "Stopped.")
	isStopped <- true
}

func (p *Parser) EmergencyStop() {
	p.CommandChan <- *p.parse("quit")
}

func (p *Parser) parse(cmdStr string) *Command {
	cmd := strings.Split(cmdStr, " ")
	return NewCommand(collectArgs(cmd))
}

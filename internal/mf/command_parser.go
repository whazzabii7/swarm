package mf

import (
	"fmt"
	"os"
	"bufio"
	"strings"

	// "github.com/whazzabii7/swarm/internal/models" 
)

type CommandParser struct {
	// requestChan	chan models.MFRequest
	CommandChan chan Command `json:"command_chan"`
}

func NewComandParser() *CommandParser {
	return &CommandParser{
		// requestChan: requests,
		CommandChan: make(chan Command),
	}
}

func (c *CommandParser) RunShell() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("swarm> ")
		if !scanner.Scan() { break }
		input := c.parse(scanner.Text())
		switch input.Type {
		case CmdQuit:
			c.CommandChan <- *input
			c.Stop()
			return
		default:
			printCommand(input)
		}
	}
}

func printCommand(c *Command) {
	count := 1
	fmt.Printf("\nCommand: %s\n", c.Type.String())
	for _, arg := range c.Args {
		fmt.Printf("\t%d: %s %s\n", count, arg.Flag, arg.Data)
		count++
	}
}

func (c * CommandParser) Stop() {
	close(c.CommandChan)
	fmt.Println("[Commander] Stopped.")
}

func (c *CommandParser) parse(cmdStr string) *Command {
	cmd := strings.Split(cmdStr, " ")
	return NewCommand(collectArgs(cmd))
}

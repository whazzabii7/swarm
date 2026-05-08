package mf

import (
	"fmt"
	"os"
	"bufio"
	"strings"

	"github.com/whazzabii7/swarm/internal/models" 
)


type Arg struct {
	Flag string  `json:"flag"`
	Data string  `json:"data"`
}

type Command struct {
	Commandtype CommandType `json:"command_type"`
	Args []Arg              `json:"args"`
}

type CommandParser struct {
	// requestChan	chan models.MFRequest
	CommandRequest chan Command
}

func NewComandParser(requests chan models.MFRequest) *CommandParser {
	return &CommandParser{
		// requestChan: requests,
		CommandRequest: make(chan Command),
	}
}

func (c *CommandParser) RunShell() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("swarm> ")
		if !scanner.Scan() { break }
		input := c.parse(scanner.Text())
		switch input.Commandtype {
		case CmdQuit:
			c.CommandRequest <- input
			c.Stop()
			return
		default:
			fmt.Printf("Command was %d", input.Commandtype)
		}
	}
}

func (c * CommandParser) Stop() {
	close(c.CommandRequest)
	fmt.Println("[Commander] Stopped.")
}

func (c *CommandParser) parse(cmdStr string) Command {
	cmd := strings.Split(cmdStr, " ")
	cmdType := stringToCommandType(cmd[0])
	return Command{
		Commandtype: cmdType,	
		Args: make([]Arg, 100),
	}
}

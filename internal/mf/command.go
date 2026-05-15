package mf

import(
	"fmt"
)

type Flag [2]string

var(
	FlagPath    = Flag { "--path", "-p" }
	FlagSource  = Flag { "--source", "-s" }
	FlagVerbose = Flag { "--verbose", "-v" }
	FlagID 		= Flag { "--id", "-i" }
	FlagAlias   = Flag { "--alias" }
	FlagPID     = Flag { "--pid" }
)


type Arg struct {
	Data []string  `json:"data"`
	IsSet bool     `json:"is_set"`
}

func NewArg(d []string, s bool) *Arg {
	return &Arg{
		Data: d,
		IsSet: s,
	}
}

func collectArgs(argStrings []string) (CommandType, map[Flag]Arg) {
	if len(argStrings) == 0 {
		return CmdPrintHelp, nil
	}

	cmdType := StringToCommandType(argStrings[0])
	
	if len(argStrings) == 1 {
		return cmdType, nil
	}

	rawArgs := argStrings[1:]
	argsMap := make(map[Flag]Arg)
	var buffer []string

	for _, current := range rawArgs {
		if len(current) > 0 && current[0] == '-' && len(buffer) > 0 {
			flagKey, arg, isInvalid := validateArg(cmdType, buffer)
			if isInvalid {
				return CmdPrintHelp, nil 
			}
			argsMap[flagKey] = *arg
			buffer = buffer[:0] 
		}
		
		buffer = append(buffer, current)
	}

	if len(buffer) > 0 {
		flagKey, arg, isInvalid := validateArg(cmdType, buffer)
		if isInvalid {
			return CmdPrintHelp, nil
		}
		argsMap[flagKey] = *arg
	}

	return cmdType, argsMap
}

func validateArg(t CommandType, buffer []string) (Flag, *Arg, bool) {
	userInput := buffer[0]
	data := buffer[1:]

	cmdFlags := map[CommandType]map[Flag]int{
		CmdListBlueprints: { FlagVerbose:0 },
		CmdListInstances: { FlagVerbose:0 },
		CmdListTasks: { FlagVerbose:0 },
		CmdListenToBot: { FlagVerbose:0 },
		CmdLoadTask: {},
		CmdPrintDBTable: { FlagVerbose:0 },
		CmdScanBotDir: { FlagPath:1 },
		CmdShowOutput: { FlagVerbose:0 },
		CmdSpawnBot: { FlagAlias:1 },
		CmdStopBot: { FlagPID:1 },
		CmdQuit: {},
		CmdPrintHelp: { FlagVerbose:1 },
	}

	if allowedFlags, ok := cmdFlags[t]; ok {
		for flagKey, paramLen := range allowedFlags {
			if (flagKey[0] != "" && flagKey[0] == userInput) || (flagKey[1] != "" && flagKey[1] == userInput) {
				
				if len(data) > paramLen {
					return Flag{}, NewArg([]string{fmt.Sprintf("Too many Arguments for %s %s", t.String(), userInput)}, false), true
				}

				if paramLen == 0 {
					return flagKey, NewArg([]string{}, true), false
				}

				return flagKey, NewArg(data, false), false
			}
		}
	}
	
	return Flag{}, NewArg([]string{fmt.Sprintf("Flag %s for Command %s not found!", userInput, t.String())}, false), true
}

type Command struct {
	Type CommandType `json:"command_type"`
	Args map[Flag]Arg       `json:"args"`
}

func NewCommand(t CommandType,args map[Flag]Arg) *Command {
	return &Command{
		Type: t,
		Args: args,
	}
}

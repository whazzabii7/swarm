package mf

import(
	"fmt"
)

type Flag [2]string

var(
	FlagPath = Flag  { "--path", "-p" }
	FlagSource = Flag { "--source", "-s" }
	FlagVerbose = Flag  { "--verbose", "-v" }
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
	// Initialisierung der Map
	argsMap := make(map[Flag]Arg)
	var buffer []string

	for _, current := range rawArgs {
		// Wenn ein neues Flag beginnt (und der Buffer nicht leer ist), verarbeite das vorherige
		if len(current) > 0 && current[0] == '-' && len(buffer) > 0 {
			// validateArg sollte nun (Flag, Arg, error) zurückgeben
			flagKey, arg, isInvalid := validateArg(cmdType, buffer)
			if isInvalid {
				// Bei Fehlern geben wir Hilfe zurück
				return CmdPrintHelp, nil 
			}
			argsMap[flagKey] = *arg
			buffer = buffer[:0] 
		}
		
		buffer = append(buffer, current)
	}

	// Den letzten Rest im Buffer verarbeiten
	if len(buffer) > 0 {
		flagKey, arg, isInvalid := validateArg(cmdType, buffer)
		if isInvalid {
			return CmdPrintHelp, nil
		}
		argsMap[flagKey] = *arg
	}

	return cmdType, argsMap
}

// Rückgabe jetzt (Flag, *Arg, bool), damit collectArgs den Map-Key kennt
func validateArg(t CommandType, buffer []string) (Flag, *Arg, bool) {
	userInput := buffer[0]
	data := buffer[1:]

	cmdFlags := map[CommandType]map[Flag]int{
		CmdListBlueprints: {},
		CmdListInstances: {},
		CmdListTasks: {},
		CmdListenToBot: {},
		CmdLoadTask: {},
		CmdPrintDBTable: {},
		CmdScanBotDir: {},
		CmdShowOutput: {},
		CmdSpawnBot: {},
		CmdStopBot: {},
		CmdQuit: {},
		CmdPrintHelp: {},
	}

	// Wir müssen prüfen, ob der userInput zu einem der erlaubten Flags gehört
	if allowedFlags, ok := cmdFlags[t]; ok {
		for flagKey, paramLen := range allowedFlags {
			// Prüfen, ob userInput im Flag-Array [2]string enthalten ist
			if (flagKey[0] != "" && flagKey[0] == userInput) || (flagKey[1] != "" && flagKey[1] == userInput) {
				
				// Deine Validierung der Argument-Länge
				if len(data) > paramLen {
					return Flag{}, NewArg([]string{fmt.Sprintf("Too many Arguments for %s %s", t.String(), userInput)}, false), true
				}

				// Deine Logik für paramLen == 0 (Bool-Verhalten)
				if paramLen == 0 {
					return flagKey, NewArg([]string{}, true), false
				}

				// Erfolg: Wir geben das flagKey-Array für die Map zurück
				return flagKey, NewArg(data, false), false
			}
		}
	}
	
	// Deine Standard-Fehlermeldung
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

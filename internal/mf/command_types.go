package mf

type CommandType int

const (
	CmdListBlueprints CommandType = iota
	CmdListInstances
	CmdListTasks
	CmdSpawnBot
	CmdStopBot
	CmdScanBotDir
	CmdLoadTask
	CmdListenToBot
	CmdShowOutput
	CmdPrintDBTable
	CmdPrintHelp
	CmdQuit
)

var cmdNames = map[CommandType]string{
	CmdListBlueprints: "list-bp",
	CmdListInstances:  "list-bots",
	CmdListTasks:      "list-tasks",
	CmdSpawnBot:       "spawn",
	CmdStopBot:        "stop",
	CmdScanBotDir:     "scan",
	CmdLoadTask:       "load",
	CmdListenToBot:    "listen",
	CmdShowOutput:     "show",
	CmdPrintDBTable:   "print-db",
	CmdPrintHelp:      "help",
	CmdQuit:           "quit",
}

var cmdAliases = map[string]CommandType{
	"q":  CmdQuit,
	":q": CmdQuit,
	"h":  CmdPrintHelp,
	"?":  CmdPrintHelp,
}

func init() {
	for cmd, name := range cmdNames {
		cmdAliases[name] = cmd
	}
}

func (c CommandType) String() string {
	if name, ok := cmdNames[c]; ok {
		return name
	}
	return "unknown"
}

func StringToCommandType(cmdStr string) CommandType {
	if cmd, ok := cmdAliases[cmdStr]; ok {
		return cmd
	}
	return CmdPrintHelp
}

package command

type CommandType int

const (
	ListBlueprints CommandType = iota
	ListInstances
	ListTasks
	SpawnBot
	StopBot
	ScanBotDir
	LoadTask
	ListenToBot
	ShowOutput
	PrintDBTable
	PrintHelp
	Quit
)

var cmdNames = map[CommandType]string{
	ListBlueprints: "list-bp",
	ListInstances:  "list-bots",
	ListTasks:      "list-tasks",
	SpawnBot:       "spawn",
	StopBot:        "stop",
	ScanBotDir:     "scan",
	LoadTask:       "load",
	ListenToBot:    "listen",
	ShowOutput:     "show",
	PrintDBTable:   "print-db",
	PrintHelp:      "help",
	Quit:           "quit",
}

var cmdAliases = map[string]CommandType{
	"q":  Quit,
	":q": Quit,
	"h":  PrintHelp,
	"?":  PrintHelp,
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
	return PrintHelp
}

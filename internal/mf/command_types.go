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

func stringToCommandType(cmdStr string) CommandType {
	switch cmdStr {
	case "quit", "q", ":q":
		return CmdQuit
	case "list-bp":
		return CmdListBlueprints
	case "list-bots":
		return CmdListInstances
	case "list-tasks":
		return CmdListTasks
	case "spawn":
		return CmdSpawnBot
	case "stop":
		return CmdStopBot
	case "load":
		return CmdLoadTask
	case "listen":
		return CmdListenToBot
	case "show":
		return CmdShowOutput
	case "print-db":
		return CmdPrintDBTable
	default:
		return CmdPrintHelp
	}
}

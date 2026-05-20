package mf

import(

	"github.com/whazzabii7/swarm/internal/db"
	"github.com/whazzabii7/swarm/internal/bot"
	"github.com/whazzabii7/swarm/internal/ui"
	"github.com/whazzabii7/swarm/internal/models"
)

func (m *Mainframe) execSpawnBot(cmd Command) {
	if arg, ok := cmd.Args[FlagAlias]; ok {
		var botBlueprint models.BotBlueprint
		if data, ok := m.blueprints[arg.Data[0]]; ok {
			botBlueprint = data
		} else {
			responseCh := make(chan models.Response)
			m.guardian.Submit(db.DBGetBlueprint, arg.Data[0], responseCh)
			response := <-responseCh
			if response.Err != nil { panic(response.Err) }
			m.blueprints[arg.Data[0]] = response.Payload.(models.BotBlueprint)
			botBlueprint = response.Payload.(models.BotBlueprint)
		}
		go m.manager.Submit(bot.BRStartBot, botBlueprint, nil)
	}
}

func (m *Mainframe) execListBlueprints(cmd Command) {}
func (m *Mainframe) execListInstances(cmd Command) {}
func (m *Mainframe) execListTasks(cmd Command) {}
func (m *Mainframe) execStopBot(cmd Command) {}
func (m *Mainframe) execScanBotDir(cmd Command) {}
func (m *Mainframe) execLoadTask(cmd Command) {}
func (m *Mainframe) execListenToBot(cmd Command) {}
func (m *Mainframe) execShowOutput(cmd Command) {}
func (m *Mainframe) execPrintDBTable(cmd Command) {}
func (m *Mainframe) execPrintHelp(cmd Command) {
	if msg, ok := cmd.Args[FlagVerbose]; ok {
		ui.Logf(ui.LevelInfo, "Help", "%s\n", msg.Data[0])
	} else {
		ui.Log(ui.LevelInfo, "Help", "Not a command!!!")
	}
}


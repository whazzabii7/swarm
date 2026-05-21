package mf

import(

	"github.com/whazzabii7/swarm/internal/db"
	"github.com/whazzabii7/swarm/internal/mf/command"
	"github.com/whazzabii7/swarm/internal/bot"
	"github.com/whazzabii7/swarm/internal/ui"
	"github.com/whazzabii7/swarm/internal/models"
)

func (m *Mainframe) execSpawnBot(cmd command.Command) {
	if arg, ok := cmd.Args[command.FlagAlias]; ok {
		var botBlueprint models.BotBlueprint
		if data, ok := m.blueprints[arg.Data[0]]; ok {
			botBlueprint = data
		} else {
			responseCh := make(chan models.Response)
			m.guardian.Submit(db.DBGetBlueprint, arg.Data[0], responseCh)
			response := <-responseCh
			if response.Err != nil { panic(response.Err) }
			botBlueprint, ok = response.Payload.GetBluePrint()
			if ok {
				m.blueprints[arg.Data[0]] = botBlueprint
			}
		}
		go m.manager.Submit(bot.BRStartBot, botBlueprint, nil)
	}
}

func (m *Mainframe) execListBlueprints(cmd command.Command) {}
func (m *Mainframe) execListInstances(cmd command.Command) {}
func (m *Mainframe) execListTasks(cmd command.Command) {}
func (m *Mainframe) execStopBot(cmd command.Command) {}

func (m *Mainframe) execScanBotDir(cmd command.Command) {
	if arg, ok := cmd.Args[command.FlagPath]; ok {
		path := arg.Data[0]
		m.scanBotDir(path)
	}
}

func (m *Mainframe) scanBotDir(path string) {
	type syncArgs struct {
		path string
		ram map[string]models.BotBlueprint
	}
	payload := syncArgs{ path, m.blueprints }
	response := make(chan models.Response)
	go m.manager.Submit(bot.BRSyncBlueprints, payload, response)
	<-response

	go m.guardian.Submit(db.DBCheckBlueprints, m.blueprints, nil)
}

func (m *Mainframe) execLoadTask(cmd command.Command) {}
func (m *Mainframe) execListenToBot(cmd command.Command) {}
func (m *Mainframe) execShowOutput(cmd command.Command) {}
func (m *Mainframe) execPrintDBTable(cmd command.Command) {}

func (m *Mainframe) execPrintHelp(cmd command.Command) {
	if msg, ok := cmd.Args[command.FlagVerbose]; ok {
		ui.Logf(ui.LevelInfo, "Help", "%s\n", msg.Data[0])
	} else {
		ui.Log(ui.LevelInfo, "Help", "Not a command!!!")
	}
}


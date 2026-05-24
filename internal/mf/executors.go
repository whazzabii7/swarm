package mf

import (
	"fmt"

	"github.com/whazzabii7/swarm/internal/bot"
	"github.com/whazzabii7/swarm/internal/db"
	"github.com/whazzabii7/swarm/internal/mf/command"
	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/ui"
)

func (m *Mainframe) execSpawnBot(cmd command.Command) {
	responseCh := make(chan models.Response)
	if arg, ok := cmd.Args[command.FlagAlias]; ok {
		var getBlueprint func() models.BotBlueprint
		if data, ok := m.blueprints[arg.Data[0]]; ok {
			getBlueprint = models.PreparePayload(data)
		} else {
			models.PrepareSubmit[db.DBRequest, string](m.guardian.Submit, db.DBGetBlueprint, arg.Data[0], responseCh)
			var response models.Response
			if response, ok = models.CheckResponse(responseCh); !ok {
				err := fmt.Errorf("%w: %v", ErrExecFailed, response.Err)
				models.PrepareSubmit[models.MFRequest, error](m.Submit, models.MFHandleError, err, nil)
			}
			if getBlueprint, ok = models.UnwrapPayload[func() models.BotBlueprint](response.Payload); ok {
				ui.Log(ui.LevelDebug, "DEBUG", "Hier Hängts") // <--- getBlueprint wird hier nicht beschrieben, ok ist false an dieser stelle
				m.blueprints[arg.Data[0]] = getBlueprint()    // das bedeutet UnwrapPayload gibt keine funktion zurück, der Cast schlägt fehl
			}
		}
		go m.manager.Submit(bot.BRStartBot, getBlueprint, responseCh)
		if response, ok := models.CheckResponse(responseCh); !ok {
			err := fmt.Errorf("%w: %v", ErrExecFailed, response.Err)
			models.PrepareSubmit[models.MFRequest, error](m.Submit, models.MFHandleError, err, nil)
		}
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
		return
	} else if _, ok := cmd.Args[command.FlagDefault]; ok {
		m.scanBotDir("./bots")
		return
	}
	err := fmt.Errorf("%w scanBotDir: %w", ErrExecFailed, ErrWrongArguments)
	models.PrepareSubmit[models.MFRequest, error](m.Submit, models.MFHandleError, err, nil)
}

func (m *Mainframe) scanBotDir(path string) {
	response := make(chan models.Response)
	ui.Log(ui.LevelInfo, "Mainframe", "Send request to DB-Guardian")
	go models.PrepareSubmit2[bot.BotRequest, string, []models.BotBlueprint](m.manager.Submit, bot.BRSyncBlueprints, path, m.getBlueprints(), response)
	go func() {
		var blueprintsResponse models.Response
		var ok bool
		if blueprintsResponse, ok = models.CheckResponse(response); !ok {
			err := fmt.Errorf("%w: %v", ErrExecFailed, blueprintsResponse.Err)
			models.PrepareSubmit[models.MFRequest, error](m.Submit, models.MFHandleError, err, nil)
			return
		}
		var getBlueprints func() []models.BotBlueprint
		if getBlueprints, ok = models.UnwrapPayload[func() []models.BotBlueprint](blueprintsResponse.Payload); !ok {
			err := fmt.Errorf("%w: %w", ErrExecFailed, ErrCorruptedData)
			models.PrepareSubmit[models.MFRequest, error](m.Submit, models.MFHandleError, err, nil)
			return
		}
		go m.guardian.Submit(db.DBCheckBlueprints, getBlueprints, response)
		if blueprintsResponse, ok = models.CheckResponse(response); !ok {
			err := fmt.Errorf("%w: %v", ErrExecFailed, blueprintsResponse.Err)
			models.PrepareSubmit[models.MFRequest, error](m.Submit, models.MFHandleError, err, nil)
			return
		}
		m.Submit(models.MFUpdateBlueprints, getBlueprints, nil)
	}()
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

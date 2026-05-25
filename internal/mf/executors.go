package mf

import (
	"fmt"

	"github.com/whazzabii7/swarm/internal/bot"
	"github.com/whazzabii7/swarm/internal/db"
	"github.com/whazzabii7/swarm/internal/mf/command"
	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/rpr"
	"github.com/whazzabii7/swarm/internal/ui"
)

func (m *Mainframe) execSpawnBot(cmd command.Command) {
	responseCh := make(chan *rpr.Response)
	if arg, ok := cmd.Args[command.FlagAlias]; ok {
		var getBlueprint func() models.BotBlueprint
		if data, ok := m.blueprints[arg.Data[0]]; ok {
			getBlueprint = rpr.PreparePayload(data)
		} else {
			rpr.PrepareSubmit[db.DBRequest, string](m.guardian.Submit, db.DBGetBlueprint, arg.Data[0], responseCh)
			var response *rpr.Response
			if response, ok = rpr.CheckResponse(responseCh); !ok {
				err := fmt.Errorf("%w: %v", ErrExecFailed, response.Err)
				rpr.PrepareSubmit[rpr.MFRequest, error](m.Submit, rpr.MFHandleError, err, nil)
			}
			if getBlueprint, ok = rpr.UnwrapPayload[func() models.BotBlueprint](response.Payload); ok {
				m.blueprints[arg.Data[0]] = getBlueprint()
			}
			response.Release()
		}
		go m.manager.Submit(bot.BRStartBot, getBlueprint, responseCh)
		if response, ok := rpr.CheckResponse(responseCh); !ok {
			err := fmt.Errorf("%w: %v", ErrExecFailed, response.Err)
			rpr.PrepareSubmit[rpr.MFRequest, error](m.Submit, rpr.MFHandleError, err, nil)
			response.Release()
		} else {
			response.Release()
		}
	}
}

func (m *Mainframe) execListBlueprints(cmd command.Command) {}
func (m *Mainframe) execListInstances(cmd command.Command)  {}
func (m *Mainframe) execListTasks(cmd command.Command)      {}
func (m *Mainframe) execStopBot(cmd command.Command)        {}

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
	rpr.PrepareSubmit[rpr.MFRequest, error](m.Submit, rpr.MFHandleError, err, nil)
}

func (m *Mainframe) scanBotDir(path string) {
	response := make(chan *rpr.Response)
	ui.Log(ui.LevelInfo, "Mainframe", "Send request to DB-Guardian")
	go rpr.PrepareSubmit2[bot.BotRequest, string, []models.BotBlueprint](m.manager.Submit, bot.BRSyncBlueprints, path, m.getBlueprints(), response)
	go func() {
		var blueprintsResponse *rpr.Response
		var ok bool
		var getBlueprints func() []models.BotBlueprint
		if blueprintsResponse, ok = rpr.CheckResponse(response); !ok {
			err := fmt.Errorf("%w: %w", ErrExecFailed, blueprintsResponse.Err)
			rpr.PrepareSubmit[rpr.MFRequest, error](m.Submit, rpr.MFHandleError, err, nil)
			blueprintsResponse.Release()
			return
		}
		if getBlueprints, ok = rpr.UnwrapPayload[func() []models.BotBlueprint](blueprintsResponse.Payload); !ok {
			err := fmt.Errorf("%w: %w", ErrExecFailed, ErrCorruptedData)
			rpr.PrepareSubmit[rpr.MFRequest, error](m.Submit, rpr.MFHandleError, err, nil)
			blueprintsResponse.Release()
			return
		}
		blueprintsResponse.Release()
		go m.guardian.Submit(db.DBCheckBlueprints, getBlueprints, response)
		if blueprintsResponse, ok = rpr.CheckResponse(response); !ok {
			err := fmt.Errorf("%w: %w", ErrExecFailed, blueprintsResponse.Err)
			rpr.PrepareSubmit[rpr.MFRequest, error](m.Submit, rpr.MFHandleError, err, nil)
			blueprintsResponse.Release()
			return
		}
		m.Submit(rpr.MFUpdateBlueprints, getBlueprints, nil)
		blueprintsResponse.Release()
	}()
}

func (m *Mainframe) execLoadTask(cmd command.Command)     {}
func (m *Mainframe) execListenToBot(cmd command.Command)  {}
func (m *Mainframe) execShowOutput(cmd command.Command)   {}
func (m *Mainframe) execPrintDBTable(cmd command.Command) {}

func (m *Mainframe) execPrintHelp(cmd command.Command) {
	if msg, ok := cmd.Args[command.FlagVerbose]; ok {
		ui.Logf(ui.LevelInfo, "Help", "%s\n", msg.Data[0])
	} else {
		ui.Log(ui.LevelInfo, "Help", "Not a command!!!")
	}
}

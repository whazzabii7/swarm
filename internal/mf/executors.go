package mf

import (
	"fmt"

	"github.com/whazzabii7/rpr"
	"github.com/whazzabii7/swarm/internal/bot"
	"github.com/whazzabii7/swarm/internal/db"
	"github.com/whazzabii7/swarm/internal/mf/command"
	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/ui"
)

func (m *Mainframe) executeCommand(cmd command.Command) {
	switch cmd.Type {
	case command.SpawnBot:
		m.execSpawnBot(cmd)
	case command.ListBlueprints:
		m.execListBlueprints(cmd)
	case command.ListInstances:
		m.execListInstances(cmd)
	case command.ListTasks:
		m.execListTasks(cmd)
	case command.StopBot:
		m.execStopBot(cmd)
	case command.ScanBotDir:
		m.execScanBotDir(cmd)
	case command.LoadTask:
		m.execLoadTask(cmd)
	case command.ListenToBot:
		m.execListenToBot(cmd)
	case command.ShowOutput:
		m.execShowOutput(cmd)
	case command.PrintDBTable:
		m.execPrintDBTable(cmd)
	default:
		m.execPrintHelp(cmd)
	}
}

func (m *Mainframe) execSpawnBot(cmd command.Command) {
	arg, ok := cmd.Args[command.FlagAlias]
	if !ok {
		return
	}

	responseCh := make(chan *rpr.Response)
	var blueprint models.BotBlueprint
	foundInRAM := false

	m.Submit(models.MFDataRequestRAM, rpr.Pack3(rpr.Ptr(models.PBlueprint), &arg.Data[0], rpr.Ptr(false)), responseCh)
	resRAM, ok := rpr.CheckResponse(responseCh)

	if resRAM != nil {
		if ok && resRAM.Payload != nil {
			if unpackOk := rpr.Assign(resRAM.Payload.Get1(), &blueprint); unpackOk {
				foundInRAM = true
			}
		}
		resRAM.Release()
	}

	if !foundInRAM {
		m.guardian.Submit(db.DBGetBlueprint, rpr.Pack(&arg.Data[0]), responseCh)
		resDB, ok := rpr.CheckResponse(responseCh)

		if resDB != nil {
			defer resDB.Release()
		}

		if !ok || resDB == nil {
			err := fmt.Errorf("%w: %v", ErrExecFailed, resDB.Err)
			m.Submit(models.MFHandleError, rpr.Pack(&err), nil)
			return
		}

		if unpackOk := rpr.Assign(resDB.Payload.Get1(), &blueprint); unpackOk {
			m.blueprints[arg.Data[0]] = blueprint
		} else {
			err := fmt.Errorf("%w: unpack failed", ErrExecFailed)
			m.Submit(models.MFHandleError, rpr.Pack(&err), nil)
			return
		}
	}

	go m.manager.Submit(bot.BRStartBot, rpr.Pack(&blueprint), responseCh)
	resStart, ok := rpr.CheckResponse(responseCh)
	if resStart != nil {
		defer resStart.Release()
	}

	if !ok {
		err := fmt.Errorf("%w: %v", ErrExecFailed, resStart.Err)
		m.Submit(models.MFHandleError, rpr.Pack(&err), nil)
	}
}

func (m *Mainframe) execListBlueprints(cmd command.Command) {}
func (m *Mainframe) execListInstances(cmd command.Command)  {}
func (m *Mainframe) execListTasks(cmd command.Command)      {}
func (m *Mainframe) execStopBot(cmd command.Command)        {}

func (m *Mainframe) execScanBotDir(cmd command.Command) {
	if arg, ok := cmd.Args[command.FlagPath]; ok {
		m.scanBotDir(arg.Data[0])
		return
	}

	if _, ok := cmd.Args[command.FlagDefault]; ok {
		m.scanBotDir("./bots")
		return
	}

	err := fmt.Errorf("%w scanBotDir: %w", ErrExecFailed, ErrWrongArguments)
	m.Submit(models.MFHandleError, rpr.Pack(&err), nil)
}

func (m *Mainframe) scanBotDir(path string) {
	response := make(chan *rpr.Response)
	ui.Log(ui.LevelInfo, "Mainframe", "Send request to DB-Guardian")
	go m.manager.Submit(bot.BRSyncBlueprints, rpr.Pack2(&path, rpr.Ptr(m.getBlueprints())), response)

	go m.processScanAsync(response)
}

func (m *Mainframe) processScanAsync(response chan *rpr.Response) {
	var blueprints []models.BotBlueprint

	resSync, ok := rpr.CheckResponse(response)
	if resSync != nil {
		defer resSync.Release()
	}

	if !ok {
		err := fmt.Errorf("%w: %w", ErrExecFailed, resSync.Err)
		m.Submit(models.MFHandleError, rpr.Pack(&err), nil)
		return
	}
	if unpackOk := rpr.Assign(resSync.Payload.Get1(), &blueprints); !unpackOk {
		err := fmt.Errorf("%w: %w", ErrExecFailed, ErrCorruptedData)
		m.Submit(models.MFHandleError, rpr.Pack(&err), nil)
		return
	}

	go m.guardian.Submit(db.DBCheckBlueprints, rpr.Pack(&blueprints), response)
	resCheck, ok := rpr.CheckResponse(response)
	if resCheck != nil {
		defer resCheck.Release()
	}

	if !ok {
		err := fmt.Errorf("%w: %w", ErrExecFailed, resCheck.Err)
		m.Submit(models.MFHandleError, rpr.Pack(&err), nil)
		return
	}

	m.Submit(models.MFUpdateBlueprints, rpr.Pack(&blueprints), nil)
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

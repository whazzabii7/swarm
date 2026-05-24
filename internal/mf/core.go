package mf

import (
	"context"
	"os"
	"time"
	// "encoding/json"

	"github.com/whazzabii7/swarm/internal/bot"
	"github.com/whazzabii7/swarm/internal/db"
	"github.com/whazzabii7/swarm/internal/mf/command"
	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/tasker"
	"github.com/whazzabii7/swarm/internal/ui"
)

type Mainframe struct {
	dbpath      string
	requestChan chan models.Request[models.MFRequest]

	// Submodules
	guardian *db.Guardian
	manager  *bot.BotManager
	tasker   *tasker.TaskManager
	cmder    *command.Parser
	error    *ErrorHandler

	// RAM-Memory (State)
	blueprints map[string]models.BotBlueprint // Key: Alias
	instances  map[int]models.BotInstance     // Key: PID
}

func NewMainframe() *Mainframe {
	// scan for database directory
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	// initializing request channel and Mainframe "RAM"
	// base initialization needed for submodules
	m := Mainframe{
		requestChan: make(chan models.Request[models.MFRequest], 100),
		dbpath:      "./data/swarm.db",
		blueprints:  make(map[string]models.BotBlueprint),
		instances:   make(map[int]models.BotInstance),
		error:       &ErrorHandler{},
	}

	m.guardian = db.NewGuardian()
	m.manager = bot.NewManager(m.Submit)
	m.tasker = tasker.NewTaskManager(m.requestChan)
	m.cmder = command.NewParser()
	return &m
}

func (m *Mainframe) Start(done chan bool) {
	ctx, cancel := context.WithCancel(context.Background())

	// DB initializing
	db.InitDB(ctx, m.dbpath)

	// Submodule initializing
	isStarted := make(chan bool)

	go m.guardian.Start(ctx, isStarted)
	m.wait(isStarted)
	go m.manager.Start(ctx, isStarted)
	m.wait(isStarted)
	m.scanBotDir("./bots")

	go m.cmder.RunShell()

	// main loop
	ui.Log(ui.LevelInfo, "Mainframe", "[!] running. Waiting for instructions...")
	for {
		select {
		case req := <-m.requestChan:
			go m.handleRequest(req)
		case cmd := <-m.cmder.CommandChan:
			if cmd.Type == command.Quit {
				m.shutdown(done, cancel)
				return
			}
			go m.executeCommand(cmd)
		case <-time.After(5 * time.Second):
			go m.checkHealth()
		}
	}
}

func (m *Mainframe) Submit(t models.MFRequest, data any, response chan models.Response) {
	m.requestChan <- models.NewRequest[models.MFRequest](t, data, response)
}

func (m *Mainframe) wait(cond chan bool) {
	if <-cond {
		return
	}
}

func (m *Mainframe) handleRequest(req models.Request[models.MFRequest]) {
	switch req.Type {
	case models.MFHandleError:
		if getErr, ok := models.UnwrapPayload[func() error](req.Payload); ok {
			err := getErr()
			errPolicy := m.error.Analyze(err)
			switch errPolicy.Severity {
			case SeverityFatal:
				ui.Logf(ui.LevelError, "Mainframe", "%s", errPolicy.Message)
				m.cmder.EmergencyStop()
			case SeverityRecover:
				ui.Logf(ui.LevelError, "Mainframe", "%s", errPolicy.Message)
				models.NewResponseErr(err).Submit(req.Response)
			case SeverityReport:
				ui.Logf(ui.LevelError, "Mainframe", "%s", errPolicy.Message)
			case SeverityIgnore:
				return
			}
		}
	}
}

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

func (m *Mainframe) checkHealth() {}

func (m *Mainframe) shutdown(done chan bool, cancel context.CancelFunc) {
	cancel()
	isStopped := make(chan bool)
	go m.cmder.Stop(isStopped)
	m.wait(isStopped)
	go m.tasker.Stop(isStopped)
	m.wait(isStopped)
	go m.manager.Stop(isStopped)
	m.wait(isStopped)
	go m.guardian.Stop(isStopped)
	m.wait(isStopped)
	done <- true
}

func (m *Mainframe) getBlueprints() []models.BotBlueprint {
	return make([]models.BotBlueprint, 0)
}

func (m *Mainframe) getInstances() []models.BotInstance {
	return make([]models.BotInstance, 0)
}

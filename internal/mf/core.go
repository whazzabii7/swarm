package mf

import (
	"context"
	"os"
	"time"
	// "encoding/json"

	"github.com/whazzabii7/rpr"
	"github.com/whazzabii7/swarm/internal/bot"
	"github.com/whazzabii7/swarm/internal/db"
	"github.com/whazzabii7/swarm/internal/mf/command"
	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/tasker"
	"github.com/whazzabii7/swarm/internal/ui"
)

type Request = rpr.Request[models.MFRequest]

type Mainframe struct {
	dbpath      string
	requestChan chan *Request

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
		requestChan: rpr.MakeRequestChan[models.MFRequest](100),
		dbpath:      "./data/swarm.db",
		blueprints:  make(map[string]models.BotBlueprint),
		instances:   make(map[int]models.BotInstance),
		error:       &ErrorHandler{},
	}

	m.guardian = db.NewGuardian()
	m.manager = bot.NewManager(m.Submit)
	m.tasker = tasker.NewTaskManager(m.Submit)
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

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// main loop
	ui.Log(ui.LevelInfo, "Mainframe", "[!] running. Waiting for instructions...")
	for {
		select {
		case req, ok := <-m.requestChan:
			if !ok {
				return
			}
			m.handleRequest(req)
		case cmd := <-m.cmder.CommandChan:
			if cmd.Type == command.Quit {
				m.shutdown(done, cancel)
				continue
			}
			go m.executeCommand(cmd)
		case <-ticker.C:
			go m.checkHealth()
		}
	}
}

func (m *Mainframe) Submit(t models.MFRequest, data rpr.Payload, response chan *rpr.Response) {
	rpr.NewRequest[models.MFRequest](t, data, response).Submit(m.requestChan)
}

func (m *Mainframe) wait(cond chan bool) {
	if <-cond {
		return
	}
}

func (m *Mainframe) checkHealth() {
	// ui.Log(ui.LevelWarning, "Mainframe", "Health Check!")
}

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
	close(m.requestChan)
	done <- true
}

func (m *Mainframe) getBlueprints() []models.BotBlueprint {
	return make([]models.BotBlueprint, 0)
}

func (m *Mainframe) getInstances() []models.BotInstance {
	return make([]models.BotInstance, 0)
}

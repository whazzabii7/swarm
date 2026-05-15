package mf

import (
	"context"
	// "fmt"
	"log"
	"os"
	"time"
	// "encoding/json"

	"github.com/whazzabii7/swarm/internal/db"
	"github.com/whazzabii7/swarm/internal/bot"
	"github.com/whazzabii7/swarm/internal/tasker"
	"github.com/whazzabii7/swarm/internal/models"
)


type Mainframe struct {
	dbpath string
	requestChan chan models.Request[models.MFRequest]

	// Submodules
	guardian *db.Guardian
	manager *bot.BotManager
	tasker  *tasker.TaskManager
	cmder   *CommandParser

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
		dbpath: "./data/swarm.db",
		blueprints: make(map[string]models.BotBlueprint),
		instances: make(map[int]models.BotInstance),
	}

	m.guardian = db.NewGuardian()
	m.manager = bot.NewBotManager(m.requestChan)
	m.tasker = tasker.NewTaskManager(m.requestChan)
	m.cmder = NewComandParser()
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
	go m.cmder.RunShell()

	// main loop
	log.Println("[!] Mainframe running. Waiting for instructions...")
	for {
		select {
			case req := <-m.requestChan:
				m.handleRequest(req)
			case cmd := <-m.cmder.CommandChan:
				if cmd.Type == CmdQuit {
					m.shutdown(done, cancel)
					return
				}
				m.executeCommand(cmd)
			case <-time.After(5*time.Second):
				m.checkHealth()
		}
	}
}

func (m *Mainframe) wait(cond chan bool) { 
	if <-cond {
		return
	}
}

func (m *Mainframe) handleRequest(req models.Request[models.MFRequest]) {}

func (m *Mainframe) executeCommand(cmd Command) {
	switch cmd.Type {
	case CmdSpawnBot:
	if arg, ok := cmd.Args[FlagAlias]; ok {
		var botBlueprint models.BotBlueprint
		if data, ok := m.blueprints[arg.Data[0]]; !ok {
			botBlueprint = data
		} else {
			responseCh := make(chan models.Response)
			m.guardian.Submit(db.DBGetBlueprint, arg, responseCh)
			response := <-responseCh
			if response.Err != nil { panic(response.Err) }
			m.blueprints[arg.Data[0]] = response.Payload.(models.BotBlueprint)
			botBlueprint = response.Payload.(models.BotBlueprint)
		}
		go m.manager.Submit(bot.BRStartBot, botBlueprint, nil)
	}
	default:
	}
}

func (m *Mainframe) checkHealth() {}

func (m *Mainframe) shutdown(done chan bool, cancel context.CancelFunc) {
	cancel()
	m.tasker.Stop()
	m.manager.Stop()
	m.guardian.Stop()
	done<-true
}

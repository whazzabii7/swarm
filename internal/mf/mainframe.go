package mf

import (
	_ "fmt"
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
	guardian *db.Guardian
	manager *bot.BotManager
	tasker  *tasker.TaskManager
	cmder   *CommandParser
	dbpath string

	// RAM-Memory (State)
	blueprints map[string]models.BotBlueprint // Key: Alias
	instances  map[int]models.BotInstance     // Key: PID

	requestChan chan models.MFRequest
}

func NewMainframe() *Mainframe {
	// scan for database
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	m := Mainframe{
		requestChan: make(chan models.MFRequest, 100),
		dbpath: "./data/swarm.db",
		blueprints: make(map[string]models.BotBlueprint),
		instances: make(map[int]models.BotInstance),
	}
	m.guardian = db.NewGuardian()
	m.manager = bot.NewBotManager(m.requestChan)
	m.tasker = tasker.NewTaskManager(m.requestChan)
	m.cmder = NewComandParser(m.requestChan)
	return &m
}

func (m *Mainframe) Start(done chan bool) {
	// DB initializing
	db.InitDB(m.dbpath)
	go m.guardian.Start()

	go m.manager.Start()

	go m.cmder.RunShell()

	// main loop
	log.Println("[!] Mainframe running. Waiting for instructions...")

	for {
		select {
			case req := <-m.requestChan:
				m.handleRequest(req)
			case cmd := <-m.cmder.CommandRequest:
				if cmd.Commandtype == CmdQuit {
					m.shutdown(done)
					return
				}
				m.executeCommand(cmd)
			case <-time.After(5*time.Second):
				m.checkHealth()
		}
	}
}

func (m *Mainframe) handleRequest(req models.MFRequest) {}
func (m *Mainframe) executeCommand(cmd Command) {}
func (m *Mainframe) checkHealth() {}

func (m *Mainframe) shutdown(done chan bool) {
	m.guardian.Stop()
	m.manager.Stop()
	m.tasker.Stop()
	done<-true
}

package ui

import (
	"fmt"
	"strings"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarning
	LevelError
)

type LogEvent struct {
	Level   LogLevel
	Module  string
	Message string
}

type UIManager struct {
	logChan     chan LogEvent
	prompt      string
	inputBuffer *strings.Builder 
}

var MainUI *UIManager

func InitUI() {
	MainUI = &UIManager{
		logChan:     make(chan LogEvent, 500),
		prompt:      "swarm> ",
		inputBuffer: &strings.Builder{},
	}
	go MainUI.routeAndPrintLoop()
}

func (ui *UIManager) Stop() {
	Log(LevelWarning, "UI", "closing UI...")
	close(ui.logChan)
}

func Log(level LogLevel, module, msg string) {
	if MainUI == nil {
		return
	}
	MainUI.logChan <- LogEvent{Level: level, Module: module, Message: msg}
}

func Logf(level LogLevel, module, format string, args ...any) {
	Log(level, module, fmt.Sprintf(format, args...))
}


func (ui *UIManager) routeAndPrintLoop() {
	for event := range ui.logChan {
		
		var prefix string
		switch event.Level {
		case LevelError:
			prefix = "\x1b[31m[ERROR]\x1b[0m" 
		case LevelWarning:
			prefix = "\x1b[33m[WARN]\x1b[0m" 
		default:
			prefix = fmt.Sprintf("[%s]", event.Module)
		}

		formattedMsg := fmt.Sprintf("%s %s", prefix, event.Message)

		ui.writeToTerminal(formattedMsg)
	}
}

func (ui *UIManager) writeToTerminal(msg string) {
	fmt.Print("\r\x1b[K")

	fmt.Println(msg)

	fmt.Print(ui.prompt)
}

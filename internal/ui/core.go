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

func InitUI(uiDone chan bool) {
	MainUI = &UIManager{
		logChan:     make(chan LogEvent, 500),
		prompt:      "[swarm> ",
		inputBuffer: &strings.Builder{},
	}
	go MainUI.routeAndPrintLoop(uiDone)
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

func (ui *UIManager) routeAndPrintLoop(uiDone chan bool) {
	for event := range ui.logChan {

		var prefix string
		switch event.Level {
		case LevelError:
			prefix = "\x1b[31m[ERROR]\x1b[0m"
		case LevelWarning:
			prefix = "\x1b[33m[WARN]\x1b[0m"
		case LevelDebug:
			prefix = fmt.Sprintf("\x1b[32m[%s]\x1b[0m", event.Module)
		default:
			prefix = fmt.Sprintf("[%s]", event.Module)
		}

		formattedMsg := fmt.Sprintf("%s %s", prefix, event.Message)

		ui.writeToTerminal(formattedMsg)
	}
	uiDone <- true
}

func (ui *UIManager) writeToTerminal(msg string) {
	fmt.Print("\r\x1b[K")

	fmt.Println(msg)

	fmt.Print(ui.prompt)
}

func Debug(str string, args ...any) {
	Log(LevelDebug, "DEBUG", fmt.Sprintf(str, args))
}

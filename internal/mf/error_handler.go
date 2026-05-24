package mf

import (
	"errors"
// "github.com/whazzabii7/swarm/internal/db"
// "github.com/whazzabii7/swarm/internal/bot"
// "github.com/whazzabii7/swarm/internal/mf/command"
// "github.com/whazzabii7/swarm/internal/tasker"
)

type ErrorSeverity int

const (
	SeverityIgnore ErrorSeverity = iota
	SeverityWarn
	SeverityReport
	SeverityRecover
	SeverityFatal
)

type ErrorPolicy struct {
	Severity ErrorSeverity
	Message  string
}

type ErrorHandler struct{}

func (e *ErrorHandler) Analyze(err error) ErrorPolicy {
	var svrt ErrorSeverity
	switch  {
	case errors.Is(err, ErrExecFailed):
		svrt = SeverityReport
	default:
		svrt = SeverityReport
	}
	return ErrorPolicy{
		Severity: svrt,
		Message:  err.Error(),
	}
}

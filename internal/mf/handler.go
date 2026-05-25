package mf

import (
	"fmt"
	"strconv"

	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/rpr"
	"github.com/whazzabii7/swarm/internal/ui"
)

func (m *Mainframe) handleRequest(req *Request) {
	switch req.Type {
	case rpr.MFDataRequestRAM:
		m.handleDataRequestRAM(req)
	case rpr.MFHandleError:
	}
}

func (m *Mainframe) handlerError(req *Request) {
	var getErr func() error
	var ok bool
	if getErr, ok = rpr.UnwrapPayload[func() error](req.Payload); !ok {
		req.Release()
		return
	}
	err := getErr()
	errPolicy := m.error.Analyze(err)
	switch errPolicy.Severity {
	case SeverityFatal:
		ui.Logf(ui.LevelError, "Mainframe", "%s", errPolicy.Message)
		m.cmder.EmergencyStop()
	case SeverityRecover:
		ui.Logf(ui.LevelError, "Mainframe", "%s", errPolicy.Message)
		rpr.NewResponseErr(err).Submit(req.Response)
	case SeverityReport:
		ui.Logf(ui.LevelError, "Mainframe", "%s", errPolicy.Message)
	case SeverityIgnore:
	}
	req.Release()
}

func (m *Mainframe) handleDataRequestRAM(req *Request) {
	var getArgs func() (rpr.RAMPage, string, bool)
	var ok bool
	if getArgs, ok = rpr.UnwrapPayload[func() (rpr.RAMPage, string, bool)](req.Payload); !ok {
		req.Release()
		return
	} 
	page, key, wantDBFallback  := getArgs()
	switch page {
	case rpr.PBlueprint:
		if blueprint, exists := m.blueprints[key]; exists {
			rpr.NewResponse[models.BotBlueprint](blueprint, nil).Submit(req.Response)
			req.Release()
			return
		} else if wantDBFallback {
			var getBlueprint func() models.BotBlueprint
			var response *rpr.Response
			responseCh := make(chan *rpr.Response)
			rpr.PrepareSubmit2[rpr.MFRequest, rpr.RAMPage, string](m.Submit, rpr.MFDataRequestDB, page, key, responseCh)
			if response, ok = rpr.CheckResponse(responseCh); !ok {
				err := fmt.Errorf("%w: %w", ErrFailedFetchingData, rpr.ErrNotAResponse)
				rpr.NewResponseErr(err).Submit(req.Response)
			} 
			if getBlueprint, ok = rpr.UnwrapPayload[func() models.BotBlueprint](response.Payload); !ok {
				var err error
				if response.Err != nil {
					err = fmt.Errorf("%w: %w", ErrFailedFetchingData, response.Err)
				} else {
					err = fmt.Errorf("%w: data is corrupted!", ErrFailedFetchingData)
				}
				rpr.NewResponseErr(err).Submit(req.Response)
			}
			rpr.NewResponse[models.BotBlueprint](getBlueprint(), nil).Submit(req.Response)
		}
	case rpr.PInstance:
		id, err := strconv.Atoi(key)
		if err != nil {
			return
		}
		if instance, exists := m.instances[id]; exists {
			rpr.NewResponse[models.BotInstance](instance, err).Submit(req.Response)
			req.Release()
			return
		} else if wantDBFallback {
			rpr.PrepareSubmit2[rpr.MFRequest, rpr.RAMPage, string](m.Submit, rpr.MFDataRequestDB, page, key, req.Response)
		}
	case rpr.PTask:
		// not implemented yet
	}

	req.Release()
}

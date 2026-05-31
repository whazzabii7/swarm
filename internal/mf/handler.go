package mf

import (
	"fmt"
	"strconv"

	"github.com/whazzabii7/rpr"
	"github.com/whazzabii7/swarm/internal/models"
	"github.com/whazzabii7/swarm/internal/ui"
)

func (m *Mainframe) handleRequest(req *Request) {
	switch req.Type {
	case models.MFDataRequestRAM:
		m.handleDataRequestRAM(req)
	case models.MFHandleError:
		m.handleError(req)
	}
}

func (m *Mainframe) handleError(req *Request) {
	defer req.Release() // Sofortige Absicherung

	var err error
	if ok := rpr.Assign(req.Payload.Get1(), &err); !ok {
		return
	}

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
}

func (m *Mainframe) handleDataRequestRAM(req *Request) {
	defer req.Release()

	var page models.RAMPage
	var key string
	var wantDBFallback bool

	pay1, pay2, pay3 := req.Payload.Get3()
	ok1 := rpr.Assign(pay1, &page)
	ok2 := rpr.Assign(pay2, &key)
	ok3 := rpr.Assign(pay3, &wantDBFallback)

	if !ok1 || !ok2 || !ok3 {
		err := fmt.Errorf("%w: %w", ErrFailedFetchingData, rpr.ErrUnpackPayloadFail)
		rpr.NewResponseErr(err).Submit(req.Response)
		return
	}

	switch page {
	case models.PBlueprint:
		m.serveRAMBlueprint(req, key, wantDBFallback)
	case models.PInstance:
		m.serveRAMInstance(req, key, wantDBFallback)
	case models.PTask:
		// not implemented yet
		rpr.NewResponse(nil, nil).Submit(req.Response)
	default:
		rpr.NewResponse(nil, nil).Submit(req.Response)
	}
}

func (m *Mainframe) serveRAMBlueprint(req *Request, key string, wantDBFallback bool) {
	if blueprint, exists := m.blueprints[key]; exists {
		rpr.NewResponse(rpr.Pack(&blueprint), nil).Submit(req.Response)
		return
	}

	if !wantDBFallback {
		rpr.NewResponse(nil, nil).Submit(req.Response)
		return
	}

	var blueprint models.BotBlueprint
	page := models.PBlueprint
	responseCh := make(chan *rpr.Response)

	m.Submit(models.MFDataRequestDB, rpr.Pack2(&page, &key), responseCh)
	resDB, ok := rpr.CheckResponse(responseCh)

	if resDB != nil {
		defer resDB.Release()
	}

	if !ok || resDB == nil {
		err := fmt.Errorf("%w: %w", ErrFailedFetchingData, rpr.ErrNotAResponse)
		rpr.NewResponseErr(err).Submit(req.Response)
		return
	}

	if unpackOk := rpr.Assign(resDB.Payload.Get1(), &blueprint); !unpackOk {
		err := fmt.Errorf("%w: %w", ErrFailedFetchingData, resDB.Err)
		if resDB.Err == nil {
			err = fmt.Errorf("%w: data is corrupted!", ErrFailedFetchingData)
		}
		rpr.NewResponseErr(err).Submit(req.Response)
		return
	}

	rpr.NewResponse(rpr.Pack(&blueprint), nil).Submit(req.Response)
}

func (m *Mainframe) serveRAMInstance(req *Request, key string, wantDBFallback bool) {
	id, err := strconv.Atoi(key)
	if err != nil {
		rpr.NewResponseErr(fmt.Errorf("%w: %w", ErrFailedFetchingData, err)).Submit(req.Response)
		return
	}

	if instance, exists := m.instances[id]; exists {
		rpr.NewResponse(rpr.Pack(&instance), err).Submit(req.Response)
		return
	}

	if !wantDBFallback {
		rpr.NewResponse(nil, nil).Submit(req.Response)
		return
	}

	// Async DB Fetch
	page := models.PInstance
	responseCh := make(chan *rpr.Response)
	m.Submit(models.MFDataRequestDB, rpr.Pack2(&page, &key), responseCh)

	go func() {
		var instance models.BotInstance
		resDB, ok := rpr.CheckResponse(responseCh)
		if resDB != nil {
			defer resDB.Release()
		}

		if !ok {
			err := fmt.Errorf("%w: %w", ErrFailedFetchingData, resDB.Err)
			rpr.NewResponseErr(err).Submit(req.Response)
			return
		}

		rpr.Assign(req.Payload.Get1(), &instance)
		rpr.NewResponse(rpr.Pack(&instance), nil).Submit(req.Response)
	}()
}

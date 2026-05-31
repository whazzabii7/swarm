package models

import "github.com/whazzabii7/rpr"

type MFSubmit = func(MFRequest, rpr.Payload, chan *rpr.Response)

type MFRequest rpr.RequestType

const (
	MFDataRequestRAM MFRequest = iota
	MFDataRequestDB
	MFHandleError
	MFUpdateBlueprints
	MFDebug
)

type RAMPage int

const (
	PBlueprint RAMPage = iota
	PInstance
	PTask
)

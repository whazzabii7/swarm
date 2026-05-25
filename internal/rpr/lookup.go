package rpr

type MFSubmit = func(MFRequest, any, chan *Response)

type MFRequest RequestType
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

package rpr

import "unsafe"

type RequestType int
type RequestConstraint interface{ ~int }

type MFRequest RequestType
type MFSubmit = func(MFRequest, any, chan *Response)

const (
	MFDataRequest MFRequest = iota
	MFHandleError
	MFUpdateBlueprints
	MFDebug
)

type rawRequest struct {
	Type     int
	Payload  Payload
	Response chan *Response
}

type Request[T RequestConstraint] struct {
	Type     T              `json:"type"`
	Payload  Payload        `json:"payload"`
	Response chan *Response `json:"response"`
}

func NewRequest[T RequestConstraint](reqType T, payload any, response chan *Response) *Request[T] {
	raw := requestPool.Get().(*rawRequest)

	req := (*Request[T])(unsafe.Pointer(raw))

	req.Type = reqType
	req.Payload = WrapPayload(payload)
	req.Response = response

	return req
}

func (r *Request[T]) Release() {
	if r == nil {
		return
	}

	r.Payload = Payload{}
	r.Response = nil

	requestPool.Put((*rawRequest)(unsafe.Pointer(r)))
}

package rpr

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

type Request[T RequestConstraint] struct {
	Type     T              `json:"type"`
	Payload  Payload        `json:"payload"`
	Response chan *Response `json:"response"`
}

func NewRequest[T RequestConstraint](typ T, payload any, response chan *Response) *Request[T] {
	return &Request[T]{
		Type:     typ,
		Payload:  WrapPayload(payload),
		Response: response,
	}
}

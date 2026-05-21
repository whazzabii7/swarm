package models

type RequestType int
type RequestConstraint interface { ~int }

type MFRequest RequestType

const (
	MFDataRequest MFRequest = iota
)

type Request[T RequestConstraint] struct {
	Type        T             `json:"type"`
	Payload     payload       `json:"payload"`
	Response    chan Response `json:"response"`
}

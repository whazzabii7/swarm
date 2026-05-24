package bot

import (
	"context"

	"github.com/whazzabii7/swarm/internal/rpr"
)

type RequestListener struct {
	requestChan  chan rpr.Request[BotRequest]
	listenerChan chan ListenerMessage
}

func NewRequestListener(request chan rpr.Request[BotRequest], listen chan ListenerMessage) *RequestListener {
	return &RequestListener{
		requestChan:  request,
		listenerChan: listen,
	}
}

func (r *RequestListener) Start(ctx context.Context, isStarted chan bool) {
	isStarted <- true
	for req := range r.requestChan {
		r.listenerChan <- ListenerMessage{
			source:      ListenToMFRequest,
			requestType: req.Type,
			payload:     req.Payload,
			response:    req.Response,
		}
	}
}

func (r *RequestListener) Stop() {
	close(r.requestChan)
}

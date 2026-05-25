package bot

import (
	"context"

	"github.com/whazzabii7/swarm/internal/rpr"
)

type BotMessage struct{}

type BotListener struct {
	botMessageChan chan BotMessage
	listenerChan   chan *rpr.Request[BotRequest]
}

func NewBotListener(listen chan *rpr.Request[BotRequest]) *BotListener {
	return &BotListener{
		botMessageChan: make(chan BotMessage, 100),
		listenerChan:   listen,
	}
}

func (b *BotListener) Start(ctx context.Context, isStarted chan bool) {
	isStarted <- true
	for msg := range b.botMessageChan {
		switch msg {
		default:
		}
	}
}

func (b *BotListener) Stop() {
	close(b.botMessageChan)
}

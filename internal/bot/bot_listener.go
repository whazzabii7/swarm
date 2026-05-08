package bot

type BotMessage struct {}

type BotListener struct {
	botMessageChan chan BotMessage
	listenerChan chan ListenerMessage
}

func NewBotListener(listen chan ListenerMessage) *BotListener {
	return &BotListener{
		botMessageChan: make(chan BotMessage),
		listenerChan: listen,
	}
}

func (b *BotListener) Start() {
	for msg := range b.botMessageChan {
		switch msg {
		default:
		}
	}
}

func (b *BotListener) Stop() {
	close(b.botMessageChan)
}

package bot

type RequestListener struct {
	requestChan chan BotRequest
	listenerChan chan ListenerMessage
}

func NewRequestListener(request chan BotRequest, listen chan ListenerMessage) *RequestListener {
	return &RequestListener{
		requestChan: request,
		listenerChan: listen,
	}
}

func (r *RequestListener) Start() {
	for req := range r.requestChan {
		switch req {
		default:
		}
	}
}

func (r *RequestListener) Stop() {
	close(r.requestChan)
}

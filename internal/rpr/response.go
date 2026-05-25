package rpr

type Response struct {
	Payload Payload `json:"payload"`
	Err     error   `json:"err"`
}

func NewResponse[T any](payload T, err error) *Response {
	res := responsePool.Get().(*Response)

	res.Payload = WrapPayload(PreparePayload(payload))
	res.Err = err

	return res
}

func NewResponse2[T1 any, T2 any](payload1 T1, payload2 T2, err error) *Response {
	res := responsePool.Get().(*Response)

	res.Payload = WrapPayload(PreparePayload2(payload1, payload2))
	res.Err = err

	return res
}

func NewResponse3[T1 any, T2 any, T3 any](payload1 T1, payload2 T2, payload3 T3, err error) *Response {
	res := responsePool.Get().(*Response)

	res.Payload = WrapPayload(PreparePayload3(payload1, payload2, payload3))
	res.Err = err

	return res
}

func NewResponseErr(err error) *Response {
	res := responsePool.Get().(*Response)

	res.Payload = WrapPayload(PreparePayload[any](nil))
	res.Err = err

	return res
}

func (r *Response) Release() {
	if r == nil {
		return
	}

	r.Payload = Payload{}
	r.Err = nil

	responsePool.Put(r)
}

func (r *Response) Submit(responseCh chan *Response) {
	responseCh <- r
}

func CheckResponse(response <-chan *Response) (*Response, bool) {
	responseData, ok := <-response
	if !ok {
		responseData = NewResponseErr(ErrNotAResponse)
	}
	if responseData.Err != nil {
		return responseData, false
	}
	return responseData, true
}

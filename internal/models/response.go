package models

type Response struct {
	Payload Payload `json:"payload"`
	Err     error   `json:"err"`
}

func NewResponse[T any](payload T, err error) *Response {
	return &Response{
		Payload: WrapPayload(payload),
		Err:     err,
	}
}

func NewResponse2[T1 any, T2 any](payload1 T1, payload2 T2, err error) *Response {
	return &Response{
		Payload: WrapPayload(PreparePayload2[T1, T2](payload1, payload2)),
		Err:     err,
	}
}

func NewResponse3[T1 any, T2 any, T3 any](payload1 T1, payload2 T2, payload3 T3, err error) *Response {
	return &Response{
		Payload: WrapPayload(PreparePayload3[T1, T2, T3](payload1, payload2, payload3)),
		Err:     err,
	}
}

func NewResponseErr(err error) *Response {
	return &Response{
		Payload: WrapPayload(PreparePayload[any](nil)),
		Err: err,
	}
}

func (r *Response) Submit(responseCh chan Response) {
	responseCh <- *r
}

func CheckResponse(response <-chan Response) (Response, bool) {
		responseData, ok := <-response
		if !ok {
			responseData = *NewResponseErr(ErrNotAResponse)
		}
		if responseData.Err != nil {
			return responseData, false
		}
		return responseData, true
}

package models

type Payload struct {
	data any
}

func WrapPayload(load any) Payload {
	return Payload{data: load}
}

func UnwrapPayload[T any](p Payload) (T, bool) {
	fn, ok := p.data.(T)
	return fn, ok
}

// Prepare Payload for Wrapping with 1 Data Value
func PreparePayload[T any](returnValue T) func() T {
	return func() T {
		return returnValue
	}
}

// Prepare Payload for Wrapping with 2 Data Value
func PreparePayload2[T1 any, T2 any](returnValue1 T1, returnValue2 T2) func() (T1, T2) {
	return func() (T1, T2) {
		return returnValue1, returnValue2
	}
}

// Prepare Payload for Wrapping with 3 Data Value
func PreparePayload3[T1 any, T2 any, T3 any](returnValue1 T1, returnValue2 T2, returnValue3 T3) func() (T1, T2, T3) {
	return func() (T1, T2, T3) {
		return returnValue1, returnValue2, returnValue3
	}
}

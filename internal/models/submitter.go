package models

func PrepareSubmit[R RequestConstraint, T any](
	submitFn func(R, any, chan Response),
	reqType R,
	payload T,
	response chan Response,
) {
	submitFn(reqType, PreparePayload[T](payload), response)
}

func PrepareSubmit2[R RequestConstraint, T1 any, T2 any](
	submitFn func(R, any, chan Response),
	reqType R,
	payload1 T1,
	payload2 T2,
	response chan Response,
) {
	submitFn(reqType, PreparePayload2[T1, T2](payload1, payload2), response)
}
func PrepareSubmit3[R RequestConstraint, T1 any, T2 any, T3 any](
	submitFn func(R, any, chan Response),
	reqType R,
	payload1 T1,
	payload2 T2,
	payload3 T3,
	response chan Response,
) {
	submitFn(reqType, PreparePayload3[T1, T2, T3](payload1, payload2, payload3), response)
}

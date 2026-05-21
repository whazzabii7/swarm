package models

type AllowedPayloads interface {
	string | BotBlueprint | BotInstance | []BotBlueprint
}

type payload struct {
	data any
}

func Payload(load any) payload {
	return payload{ data: load }
}

// only debuging, will not be in final build!
func (p *payload) Unwrap() any {
	return p.data
}

func castPayloadType[T AllowedPayloads](rPayload any) (T, bool) {
	value, ok := rPayload.(T)
	return value, ok
}

func (p *payload) GetBluePrint() (BotBlueprint, bool) {
	value, ok := castPayloadType[BotBlueprint](p.data)
	return value, ok
}

func (p *payload) GetBluePrints() ([]BotBlueprint, bool) {
	values, ok := castPayloadType[[]BotBlueprint](p.data)
	return values, ok
}

func (p *payload) GetString() (string, bool) {
	value, ok := castPayloadType[string](p.data)
	return value, ok
}

func (p *payload) GetInstance() (BotInstance, bool) {
	value, ok := castPayloadType[BotInstance](p.data)
	return value, ok
}

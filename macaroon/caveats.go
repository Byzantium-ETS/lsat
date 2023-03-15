package macaroon

type Caveat struct {
	Key   string
	Value string
}

func NewCaveat(Key string, Value string) Caveat {
	return Caveat{Key, Value}
}

func (caveat Caveat) ToString() string {
	return caveat.Key + "=" + caveat.Value
}

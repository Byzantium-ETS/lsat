package auth

const (
	Timeout = "timeout"
)

type Caveat struct {
	Name  string
	Value string
}

func NewCaveat(Name string, Value string) Caveat {
	return Caveat{Name, Value}
}

func (caveat Caveat) ToString() string {
	return caveat.Name + ":" + caveat.Value
}

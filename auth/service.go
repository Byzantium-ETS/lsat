package auth

type Tier = int8

// Un service est un caveat particulier
// C'est le premier caveat appliqué
type Service struct {
	Name  string
	Tier  Tier
	Price int64
}

func (service *Service) ToString() string {
	return "service=" + service.Name + ":" + string(service.Tier)
}

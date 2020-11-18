package repository

type Domain struct {
	Uuid string
	Name string
}

type IDatabase interface {
	SelectDomain(uuid string) (Domain, error)
	GetDomainsList() ([]Domain, error)
	InsertDomain(uuid, name string) error
	DeleteDomain(uuid string) error
}

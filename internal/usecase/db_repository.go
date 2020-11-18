package usecase

import "coredns_api/internal/model"

type IDatabaseRepository interface {
	GetDomain(uuid model.Uuid) (*model.Domain, error)
	GetDomainsList() ([]*model.Domain, error)
	AddDomain(domain *model.Domain) error
	DeleteDomain(domain *model.Domain) error
}

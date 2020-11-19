package repository

import (
	"coredns_api/internal/model"
	"coredns_api/internal/usecase"
)

type DatabaseRepository struct {
	db IDatabase
}

func NewDatabaseRepository(db IDatabase) usecase.IDatabaseRepository {
	return &DatabaseRepository{db}
}

func (d *DatabaseRepository) GetDomain(uuid model.Uuid) (*model.Domain, error) {
	domainResult, err := d.db.SelectDomain(uuid.String())
	if err != nil {
		return nil, err
	}

	var domain model.Domain
	domainUuid, err := model.NewUuid(domainResult.Uuid)
	if err != nil {
		return nil, err
	}

	domainName, err := model.NewDomainName(domainResult.Name)
	if err != nil {
		return nil, err
	}

	domain.Uuid = domainUuid
	domain.Name = domainName
	return &domain, nil
}

func (d *DatabaseRepository) GetDomainsList() ([]*model.Domain, error) {
	domList, err := d.db.GetDomainsList()
	if err != nil {
		return nil, err
	}

	var domainList []*model.Domain
	for _, dom := range domList {
		u, _ := model.NewUuid(dom.Uuid)
		m, _ := model.NewEmptyDomain(u, dom.Name)
		domainList = append(domainList, m)
	}
	return domainList, nil
}

func (d *DatabaseRepository) AddDomain(domain *model.Domain) error {
	uuid := domain.Uuid
	name := domain.Name
	return d.db.InsertDomain(uuid.String(), name.String())
}

func (d *DatabaseRepository) DeleteDomain(domain *model.Domain) error {
	uuid := domain.Uuid
	return d.db.DeleteDomain(uuid.String())
}

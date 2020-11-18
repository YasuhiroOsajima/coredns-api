package repository

import (
	"errors"
	"sync"

	"coredns_api/internal/model"
)

type DomainCache struct {
	sync.Mutex
	Cache map[model.DomainName]*model.Domain
}

func (d *DomainCache) Add(domain *model.Domain) {
	d.Cache[domain.Name] = domain
}

func (d *DomainCache) Get(domainName model.DomainName) (*model.Domain, error) {
	domain := d.Cache[domainName]
	if domain == nil {
		return nil, errors.New("target domain chache is not found")
	}

	return domain, nil
}

func (d *DomainCache) Delete(domain *model.Domain) {
	delete(d.Cache, domain.Name)
}

func NewDomainCache() *DomainCache {
	return &DomainCache{Cache: map[model.DomainName]*model.Domain{}}
}

var domainCache = NewDomainCache()

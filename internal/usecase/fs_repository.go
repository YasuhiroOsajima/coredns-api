package usecase

import "coredns_api/internal/model"

type IFilesystemRepository interface {
	WriteDomainFile(domain *model.Domain) error
	LoadDomainFile(domainName string) (*model.Domain, error)
	DeleteDomainFile(domain *model.Domain) error
}

package usecase

import "coredns_api/internal/model"

type IFilesystemRepository interface {
	Initialize()
	Lock()
	UnLock()
	WriteConfCache() error
	WriteDomainFile(domain *model.Domain) error
	LoadAllDomains() ([]*model.Domain, error)
	GetDomainByUuid(domainUuid model.Uuid) (*model.Domain, error)
	DeleteDomainFile(domain *model.Domain) error
}

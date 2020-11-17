package usecase

import (
	"coredns_api/internal/model"
)

type DomainInteractor struct {
	FsRepository IFilesystemRepository
}

func NewDomainInteractor(repo IFilesystemRepository) *DomainInteractor {
	return &DomainInteractor{repo}
}

func (i *DomainInteractor) Add(domain *model.Domain) error {
	err := i.FsRepository.WriteDomainFile(domain)
	if err != nil {
		return err
	}

	return nil
}

func (i *DomainInteractor) Get(domain *model.Domain) (*model.Domain, error) {
	domainName := domain.Name
	gotDomain, err := i.FsRepository.LoadDomainFile(domainName)
	if err != nil {
		return nil, err
	}

	return gotDomain, nil
}

func (i *DomainInteractor) Delete(domain *model.Domain) error {
	err := i.FsRepository.DeleteDomainFile(domain)
	if err != nil {
		return err
	}

	return nil
}

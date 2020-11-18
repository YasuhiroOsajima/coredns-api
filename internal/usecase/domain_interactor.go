package usecase

import (
	"coredns_api/internal/model"
)

type DomainInteractor struct {
	FsRepository IFilesystemRepository
	DbRepository IDatabaseRepository
}

func NewDomainInteractor(fRepo IFilesystemRepository, dRepo IDatabaseRepository) *DomainInteractor {
	return &DomainInteractor{FsRepository: fRepo, DbRepository: dRepo}
}

func (i *DomainInteractor) Add(domain *model.Domain) error {
	err := i.FsRepository.WriteDomainFile(domain)
	if err != nil {
		return err
	}

	err = i.DbRepository.AddDomain(domain)
	if err != nil {
		return err
	}

	return nil
}

func (i *DomainInteractor) Get(domainUuid model.Uuid) (*model.Domain, error) {
	targetDomain, err := i.DbRepository.GetDomain(domainUuid)
	if err != nil {
		return nil, err
	}

	gotDomainInfo, err := i.FsRepository.LoadDomainFile(targetDomain)
	if err != nil {
		return nil, err
	}

	return gotDomainInfo, nil
}

func (i *DomainInteractor) Delete(domainUuid model.Uuid) error {
	domain, err := i.DbRepository.GetDomain(domainUuid)
	if err != nil {
		return err
	}

	err = i.FsRepository.DeleteDomainFile(domain)
	if err != nil {
		return err
	}

	err = i.DbRepository.DeleteDomain(domain)
	if err != nil {
		return err
	}

	return nil
}

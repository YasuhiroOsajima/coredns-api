package usecase

import (
	"coredns_api/internal/model"
)

type DomainInteractor struct {
	fsRepository IFilesystemRepository
	dbRepository IDatabaseRepository
}

func NewDomainInteractor(fRepo IFilesystemRepository, dRepo IDatabaseRepository) *DomainInteractor {
	d := &DomainInteractor{fsRepository: fRepo, dbRepository: dRepo}

	d.fsRepository.Lock()
	defer d.fsRepository.UnLock()

	allDomainInfo, err := d.fsRepository.LoadAllDomains()
	if err != nil {
		panic(err)
	}

	for _, dom := range allDomainInfo {
		err = d.dbRepository.AddDomain(dom)
		if err != nil {
			panic(err)
		}
	}

	return d
}

func (i *DomainInteractor) Add(domain *model.Domain) error {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	err := i.fsRepository.WriteDomainFile(domain)
	if err != nil {
		return err
	}

	err = i.fsRepository.WriteConfCache()
	if err != nil {
		_ = i.fsRepository.DeleteDomainFile(domain)
		return err
	}

	err = i.dbRepository.AddDomain(domain)
	if err != nil {
		_ = i.fsRepository.DeleteDomainFile(domain)
		_ = i.fsRepository.WriteConfCache()
		return err
	}

	return nil
}

func (i *DomainInteractor) Get(domainUuid model.Uuid) (*model.Domain, error) {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	targetDomain, err := i.dbRepository.GetDomain(domainUuid)
	if err != nil {
		return nil, err
	}

	gotDomainInfo, err := i.fsRepository.LoadDomainFile(targetDomain.Name)
	if err != nil {
		return nil, err
	}

	return gotDomainInfo, nil
}

func (i *DomainInteractor) Delete(domainUuid model.Uuid) error {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	domain, err := i.dbRepository.GetDomain(domainUuid)
	if err != nil {
		return err
	}

	err = i.fsRepository.DeleteDomainFile(domain)
	if err != nil {
		return err
	}

	err = i.fsRepository.WriteConfCache()
	if err != nil {
		_ = i.fsRepository.WriteDomainFile(domain)
		return err
	}

	err = i.dbRepository.DeleteDomain(domain)
	if err != nil {
		_ = i.fsRepository.WriteDomainFile(domain)
		_ = i.fsRepository.WriteConfCache()
		return err
	}

	return nil
}

func (i *DomainInteractor) GetDomainsList() ([]*model.Domain, error) {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	return i.dbRepository.GetDomainsList()
}

package usecase

import "coredns_api/internal/model"

type DomainInteractor struct {
	fsRepository IFilesystemRepository
}

func NewDomainInteractor(fRepo IFilesystemRepository) *DomainInteractor {
	r := &DomainInteractor{fsRepository: fRepo}
	r.fsRepository.Initialize()
	return r
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

	return nil
}

func (i *DomainInteractor) Get(domainUuid model.Uuid, requestTenantUuid model.Uuid) (*model.Domain, error) {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	targetDomain, err := i.fsRepository.GetDomainByUuid(domainUuid, requestTenantUuid)
	if err != nil {
		return nil, err
	}

	return targetDomain, nil
}

func (i *DomainInteractor) Update(domainUuid model.Uuid, requestTenantUuid model.Uuid, tenantUuidList []model.Uuid) (*model.Domain, error) {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	domain, err := i.fsRepository.GetDomainByUuid(domainUuid, requestTenantUuid)
	if err != nil {
		return nil, err
	}

	err = domain.UpdateTenants(requestTenantUuid, tenantUuidList)
	if err != nil {
		return nil, err
	}

	err = i.fsRepository.WriteDomainFile(domain)
	if err != nil {
		return nil, err
	}

	err = i.fsRepository.WriteConfCache()
	if err != nil {
		_ = i.fsRepository.DeleteDomainFile(domain)
		return nil, err
	}

	return domain, nil
}

func (i *DomainInteractor) Delete(domainUuid model.Uuid, requestTenantUuid model.Uuid) error {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	domain, err := i.fsRepository.GetDomainByUuid(domainUuid, requestTenantUuid)
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

	return nil
}

func (i *DomainInteractor) GetDomainsList(requestTenantUuid model.Uuid) ([]*model.Domain, error) {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	return i.fsRepository.LoadTenantAllDomains(requestTenantUuid)
}

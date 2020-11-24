package usecase

import (
	"coredns_api/internal/model"
)

type HostInteractor struct {
	fsRepository IFilesystemRepository
}

func NewHostInteractor(fRepo IFilesystemRepository) *HostInteractor {
	return &HostInteractor{fRepo}
}

func (i *HostInteractor) Add(newHost *model.Host, domainUuid model.Uuid, requestTenantUuid model.Uuid) (*model.Domain, error) {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	gotDomain, err := i.fsRepository.GetDomainByUuid(domainUuid, requestTenantUuid)
	if err != nil {
		return nil, err
	}

	for _, h := range gotDomain.Hosts {
		if h.Name == newHost.Name {
			return nil, NewHostDuplicatedError("hostname", newHost.Name)
		}
		if h.Address == newHost.Address {
			return nil, NewHostDuplicatedError("address", newHost.Address)
		}
	}

	hosts := append(gotDomain.Hosts, newHost)
	gotDomain.Hosts = hosts

	err = i.fsRepository.WriteDomainFile(gotDomain)
	if err != nil {
		return nil, err
	}

	return gotDomain, nil
}

func (i *HostInteractor) Get(hostUuid, domainUuid model.Uuid, requestTenantUuid model.Uuid) (*model.Host, error) {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	domain, err := i.fsRepository.GetDomainByUuid(domainUuid, requestTenantUuid)
	if err != nil {
		return nil, err
	}

	for _, h := range domain.Hosts {
		if h.Uuid == hostUuid {
			return h, nil
		}
	}

	return nil, model.NewHostNotFoundError()
}

func (i *HostInteractor) GetDomain(domainUuid model.Uuid, requestTenantUuid model.Uuid) (*model.Domain, error) {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	targetDomain, err := i.fsRepository.GetDomainByUuid(domainUuid, requestTenantUuid)
	if err != nil {
		return nil, err
	}

	return targetDomain, nil
}

func (i *HostInteractor) Update(newHost *model.Host, domainUuid model.Uuid, requestTenantUuid model.Uuid) error {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	domain, err := i.fsRepository.GetDomainByUuid(domainUuid, requestTenantUuid)
	if err != nil {
		return err
	}

	var newHosts []*model.Host
	found := false
	for _, h := range domain.Hosts {
		if h.Name == newHost.Name {
			return NewHostDuplicatedError("hostname", newHost.Name)
		}
		if h.Address == newHost.Address {
			return NewHostDuplicatedError("address", newHost.Address)
		}

		if h.Uuid == newHost.Uuid {
			newHosts = append(newHosts, newHost)
			found = true
		} else {
			newHosts = append(newHosts, h)
		}
	}

	if !found {
		return model.NewHostNotFoundError()
	}

	domain.Hosts = newHosts
	return i.fsRepository.WriteDomainFile(domain)
}

func (i *HostInteractor) Delete(host *model.Host, domainUuid model.Uuid, requestTenantUuid model.Uuid) error {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	domain, err := i.fsRepository.GetDomainByUuid(domainUuid, requestTenantUuid)
	if err != nil {
		return err
	}

	var newHosts []*model.Host
	found := false
	for _, h := range domain.Hosts {
		if h.Uuid == host.Uuid {
			found = true
		} else {
			newHosts = append(newHosts, h)
		}
	}

	if !found {
		return model.NewHostNotFoundError()
	}

	domain.Hosts = newHosts
	return i.fsRepository.WriteDomainFile(domain)
}

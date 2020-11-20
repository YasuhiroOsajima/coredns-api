package usecase

import (
	"errors"

	"coredns_api/internal/model"
)

type HostInteractor struct {
	fsRepository IFilesystemRepository
	dbRepository IDatabaseRepository
}

func NewHostInteractor(fRepo IFilesystemRepository, dRepo IDatabaseRepository) *HostInteractor {
	return &HostInteractor{fRepo, dRepo}
}

func (i *HostInteractor) Add(newHost *model.Host, domainUuid model.Uuid) (*model.Domain, error) {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	gotDomain, err := i.GetDomain(domainUuid)
	if err != nil {
		return nil, err
	}

	for _, h := range gotDomain.Hosts {
		if h.Name == newHost.Name {
			return nil, errors.New("target hostname is already assigned in the domain")
		}
		if h.Address == newHost.Address {
			return nil, errors.New("target address is already assigned in the domain")
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

func (i *HostInteractor) Get(hostUuid, domainUuid model.Uuid) (*model.Host, error) {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	domain, err := i.GetDomain(domainUuid)
	if err != nil {
		return nil, err
	}

	for _, h := range domain.Hosts {
		if h.Uuid == hostUuid {
			return h, nil
		}
	}

	return nil, errors.New("target host is not found")
}

func (i *HostInteractor) GetDomain(domainUuid model.Uuid) (*model.Domain, error) {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	targetDomain, err := i.dbRepository.GetDomain(domainUuid)
	if err != nil {
		return nil, err
	}

	gotDomain, err := i.fsRepository.LoadDomainFile(targetDomain.Name)
	if err != nil {
		return nil, err
	}

	return gotDomain, nil
}

func (i *HostInteractor) Update(newHost *model.Host, domainUuid model.Uuid) error {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	domain, err := i.GetDomain(domainUuid)
	if err != nil {
		return err
	}

	var newHosts []*model.Host
	for _, h := range domain.Hosts {
		if h.Uuid == newHost.Uuid {
			newHosts = append(newHosts, newHost)
		} else {
			newHosts = append(newHosts, h)
		}
	}
	domain.Hosts = newHosts
	return i.fsRepository.WriteDomainFile(domain)
}

func (i *HostInteractor) Delete(host *model.Host, domainUuid model.Uuid) error {
	i.fsRepository.Lock()
	defer i.fsRepository.UnLock()

	domain, err := i.GetDomain(domainUuid)
	if err != nil {
		return err
	}

	var newHosts []*model.Host
	for _, h := range domain.Hosts {
		if h.Uuid != host.Uuid {
			newHosts = append(newHosts, h)
		}
	}
	domain.Hosts = newHosts
	return i.fsRepository.WriteDomainFile(domain)
}

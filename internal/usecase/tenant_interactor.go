package usecase

import "coredns_api/internal/model"

type TenantInteractor struct {
	fsRepository IFilesystemRepository
}

func NewTenantInteractor(fRepo IFilesystemRepository) *TenantInteractor {
	r := &TenantInteractor{fsRepository: fRepo}
	r.fsRepository.Initialize()
	return r
}

func (t *TenantInteractor) GetDomainList() ([]*model.Domain, error) {
	t.fsRepository.Lock()
	defer t.fsRepository.UnLock()

	return t.fsRepository.LoadAllDomains()
}

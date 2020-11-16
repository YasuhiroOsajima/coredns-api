package usecase

import "coredns_api/internal/model"

type HostInteractor struct {
	FsRepository IFilesystemRepository
}

func NewHostInteractor(repo IFilesystemRepository) *HostInteractor {
	return &HostInteractor{repo}
}

func (i *HostInteractor) Add(host model.Host) error {
	return nil
}

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

func (i *DomainInteractor) Add(domain model.Domain) error {
	return nil
}

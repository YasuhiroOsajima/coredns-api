package repository

import (
	"coredns_api/internal/model"
	"coredns_api/internal/usecase"
)

type FilesystemRepository struct {
	filesystem IFilesystem
}

func NewFileRepository(fs IFilesystem) usecase.IFilesystemRepository {
	return &FilesystemRepository{fs}
}

func (f *FilesystemRepository) WriteDomainFile(domain *model.Domain) error {
	domainCache.Lock()
	defer domainCache.Unlock()

	name := domain.Name
	fileInfo, err := domain.GetFileInfo()
	if err != nil {
		return err
	}

	err = f.filesystem.WriteTextFile(name.String(), fileInfo)

	domainCache.Add(domain)

	return err
}

func (f *FilesystemRepository) LoadDomainFile(targetDomain *model.Domain) (*model.Domain, error) {
	domainCache.Lock()
	defer domainCache.Unlock()

	domainName := targetDomain.Name
	domain, err := domainCache.Get(domainName)
	if err == nil {
		return domain, nil
	}

	fileInfo, err := f.filesystem.LoadTextFile(domainName.String())
	if err != nil {
		return nil, err
	}

	domain, err = model.NewDomain(domainName.String(), fileInfo)
	if err != nil {
		return nil, err
	}

	domainCache.Add(domain)

	return domain, nil
}

func (f *FilesystemRepository) DeleteDomainFile(domain *model.Domain) error {
	domainCache.Lock()
	defer domainCache.Unlock()

	domainCache.Delete(domain)

	return f.filesystem.DeleteFile(domain.Name.String())
}

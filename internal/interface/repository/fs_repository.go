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
	name := domain.Name
	fileInfo, err := domain.GetFileInfo()
	if err != nil {
		return err
	}

	err = f.filesystem.WriteTextFile(name, fileInfo)
	return err
}

func (f *FilesystemRepository) LoadDomainFile(domainName string) (*model.Domain, error) {
	fileInfo, err := f.filesystem.LoadTextFile(domainName)
	if err != nil {
		return nil, err
	}

	domain, err := model.NewDomain(domainName, fileInfo)
	if err != nil {
		return nil, err
	}

	return domain, nil
}

func (f *FilesystemRepository) DeleteDomainFile(domain *model.Domain) error {
	return f.filesystem.DeleteFile(domain.Name)
}

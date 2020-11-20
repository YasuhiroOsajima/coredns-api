package repository

import (
	"log"

	"coredns_api/internal/model"
	"coredns_api/internal/usecase"
)

var coreDNSConfCache *model.CoreDNSConf

type FilesystemRepository struct {
	filesystem IFilesystem
}

func NewFileRepository(fs IFilesystem) usecase.IFilesystemRepository {
	f := &FilesystemRepository{fs}

	allDomainInfo, err := f.loadAllDomainFiles()
	if err != nil {
		panic(err)
	}

	coreDNSConfCache = model.NewCoreDNSConf(allDomainInfo)
	return f
}

func (f *FilesystemRepository) WriteConfCache() error {
	if !coreDNSConfCache.IsLocked() {
		return usecase.NewIsNotLockedError()
	}

	confPath := coreDNSConfCache.ConfPath
	confInfo, err := coreDNSConfCache.GetFileInfo()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return f.filesystem.WriteTextFile(confPath, confInfo)
}

func (f *FilesystemRepository) WriteDomainFile(domain *model.Domain) error {
	if !coreDNSConfCache.IsLocked() {
		return usecase.NewIsNotLockedError()
	}

	name := domain.Name
	fileInfo, err := domain.GetFileInfo()
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = f.filesystem.WriteTextFile(name.String(), fileInfo)
	if err != nil {
		log.Fatal(err)
		return err
	}

	coreDNSConfCache.Add(domain)
	return nil
}

func (f *FilesystemRepository) LoadDomainFile(domainName model.DomainName) (*model.Domain, error) {
	if !coreDNSConfCache.IsLocked() {
		return nil, usecase.NewIsNotLockedError()
	}

	domain, err := coreDNSConfCache.Get(domainName)
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

	coreDNSConfCache.Add(domain)

	return domain, nil
}

func (f *FilesystemRepository) loadAllDomainFiles() ([]*model.Domain, error) {
	if !coreDNSConfCache.IsLocked() {
		return nil, usecase.NewIsNotLockedError()
	}

	domainFileDir := model.GetHostsDir()
	fileNameList, err := f.filesystem.GetFilenameList(domainFileDir)
	if err != nil {
		return nil, err
	}

	var domainList []*model.Domain
	for _, domainFile := range fileNameList {
		domainName, err := model.NewDomainName(domainFile)
		if err != nil {
			log.Fatal(domainFile)
			log.Fatal(err)
			return nil, err
		}

		domain, err := f.LoadDomainFile(domainName)
		if err != nil {
			log.Fatal(domainFile)
			log.Fatal(err)
			return nil, err
		}

		domainList = append(domainList, domain)
	}
	return domainList, nil
}

func (f *FilesystemRepository) DeleteDomainFile(domain *model.Domain) error {
	if !coreDNSConfCache.IsLocked() {
		return usecase.NewIsNotLockedError()
	}

	err := f.filesystem.DeleteFile(domain.Name.String())
	if err != nil {
		return err
	}

	coreDNSConfCache.Delete(domain)

	return nil
}

func (f *FilesystemRepository) Lock() {
	coreDNSConfCache.SetLocke()
}

func (f *FilesystemRepository) UnLock() {
	coreDNSConfCache.UnSetLocke()
}

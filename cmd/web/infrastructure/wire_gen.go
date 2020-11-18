// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package infrastructure

import (
	"coredns_api/internal/infrastructure"
	"coredns_api/internal/interface/repository"
	"coredns_api/internal/usecase"
	"coredns_api/pkg/interface/controllers"
)

// Injectors from wire.go:

func InitializeDomainController() *controllers.DomainController {
	iFilesystem := infrastructure.NewFilesystem()
	iFilesystemRepository := repository.NewFileRepository(iFilesystem)
	domainInteractor := usecase.NewDomainInteractor(iFilesystemRepository)
	domainController := controllers.NewDomainController(domainInteractor)
	return domainController
}

func InitializeHostController() *controllers.HostController {
	iFilesystem := infrastructure.NewFilesystem()
	iFilesystemRepository := repository.NewFileRepository(iFilesystem)
	hostInteractor := usecase.NewHostInteractor(iFilesystemRepository)
	hostController := controllers.NewHostController(hostInteractor)
	return hostController
}
//+build wireinject

package infrastructure

import (
	"github.com/google/wire"

	inf "coredns_api/internal/infrastructure"
	"coredns_api/internal/interface/repository"
	"coredns_api/internal/usecase"
	"coredns_api/pkg/interface/controllers"
)

func InitializeDomainController() *controllers.DomainController {
	wire.Build(
		controllers.NewDomainController,
		usecase.NewDomainInteractor,
		repository.NewFileRepository,
		inf.NewFilesystem,
	)
	return nil
}

func InitializeHostController() *controllers.HostController {
	wire.Build(
		controllers.NewHostController,
		usecase.NewHostInteractor,
		repository.NewFileRepository,
		inf.NewFilesystem,
	)
	return nil
}

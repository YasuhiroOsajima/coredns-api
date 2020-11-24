//+build wireinject

package infrastructure

import (
	"github.com/google/wire"

	inf "coredns_api/internal/infrastructure"
	"coredns_api/internal/interface/repository"
	"coredns_api/internal/usecase"
	"coredns_api/pkg/interface/controllers"
)

func InitializeTenantController() *controllers.TenantController {
	wire.Build(
		controllers.NewTenantController,
		usecase.NewTenantInteractor,
		repository.NewFileRepository,
		inf.NewFilesystem,
	)
	return nil
}

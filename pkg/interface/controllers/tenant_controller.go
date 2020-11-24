package controllers

import (
	"coredns_api/internal/model"
	"coredns_api/internal/usecase"
	"log"
	"net/http"
)

// Result
type TenantInfoResult struct {
	Tenants []TenantResult `json:"tenants"`
}

type TenantResult struct {
	Uuid    string   `json:"uuid"`
	Domains []string `json:"domains"`
}

// Controller
type TenantController struct {
	interactor *usecase.TenantInteractor
}

func NewTenantController(itr *usecase.TenantInteractor) *TenantController {
	return &TenantController{itr}
}

func (t *TenantController) List(c Context) {
	domains, err := t.interactor.GetDomainList()
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Print(err)
		return
	}

	var tenants []model.Uuid
	for _, d := range domains {
		for _, t := range d.Tenants {
			tenants = append(tenants, t)
		}
	}

	var result TenantInfoResult
	tenantList := make([]TenantResult, 0)
	for _, t := range tenants {
		var tenant TenantResult
		tenant.Uuid = t.String()

		for _, d := range domains {
			for _, tin := range d.Tenants {
				if t == tin {
					tdomains := tenant.Domains
					tdomains = append(tdomains, d.Uuid.String())
					tenant.Domains = tdomains
					continue
				}
			}
		}
		tenantList = append(tenantList, tenant)
	}
	result.Tenants = tenantList
	c.JSON(http.StatusOK, result)
}

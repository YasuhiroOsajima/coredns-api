package controllers

import (
	"errors"
	"log"
	"net/http"

	"coredns_api/internal/model"
	"coredns_api/internal/usecase"
)

// Request
type DomainRequest struct {
	Name    string   `json:"domain"`
	Tenants []string `json:"tenants"`
}

type DomainUpdateRequest struct {
	Tenants []string `json:"tenants"`
}

// Result
type DomainInfoResult struct {
	DomainResult
	Hosts []HostResult `json:"hosts"`
}

type DomainResult struct {
	Domain  string   `json:"domain"`
	Uuid    string   `json:"uuid"`
	Tenants []string `json:"tenants"`
}

type HostResult struct {
	Name    string `json:"hostname"`
	Address string `json:"address"`
	Uuid    string `json:"uuid"`
}

type DomainListResult struct {
	Domains []DomainResult `json:"domains"`
}

// error to return with HTTP 500
type UnAvailableHandlingError struct {
	err string
}

func NewUnAvailableHandlingError() error {
	return &UnAvailableHandlingError{err: "unhandled server side error"}
}

func (e *UnAvailableHandlingError) Error() string {
	return e.err
}

// Controller
type DomainController struct {
	interactor *usecase.DomainInteractor
}

func NewDomainController(itr *usecase.DomainInteractor) *DomainController {
	return &DomainController{itr}
}

// Add handler doc
// @Tags Domain
// @Summary Add new domain
// @Description Add new domain to coredns
// @Accept json
// @Produce json
// @Param domain body DomainRequest true "Request body parameter with json format"
// @Success 201 {object} DomainInfoResult
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/domains [post]
func (d *DomainController) Add(c Context) {
	var request DomainRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Print(err)
		return
	}

	name := request.Name
	tenantList := request.Tenants
	if len(tenantList) == 0 {
		NewError(c,
			http.StatusBadRequest,
			errors.New("accessible tenant uuid is not specified"))
		return
	}

	newDomain, err := model.NewOriginalDomain(name, tenantList)
	if err != nil {
		switch e := err.(type) {
		case *model.InvalidParameterGiven:
			NewError(c, http.StatusBadRequest, err)
		default:
			NewError(c,
				http.StatusInternalServerError,
				NewUnAvailableHandlingError())
			log.Print(e)
		}
		log.Print(err)
		return
	}

	err = d.interactor.Add(newDomain)
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Print(err)
		return
	}

	var tenants []string
	for _, t := range newDomain.Tenants {
		tenants = append(tenants, t.String())
	}

	hosts := make([]HostResult, 0)
	var result DomainInfoResult
	result.Domain = newDomain.Name.String()
	result.Uuid = newDomain.Uuid.String()
	result.Hosts = hosts
	result.Tenants = tenants
	c.JSON(http.StatusCreated, result)
}

// List handler doc
// @Tags Domain
// @Summary List domains
// @Description List domains from coredns
// @Produce json
// @Param Tenant header string true "Tenant UUID to set access control"
// @Success 200 {object} DomainListResult
// @Failure 500 {object} HTTPError
// @Router /v1/domains [get]
func (d *DomainController) List(c Context) {
	requestTenant := c.GetHeader("Tenant")
	if requestTenant == "" {
		NewError(c,
			http.StatusBadRequest,
			errors.New("tenant uuid header is not specified"))
		return
	}

	requestTenantUuid, err := model.NewUuid(requestTenant)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		return
	}

	domainList, err := d.interactor.GetDomainsList(requestTenantUuid)
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Print(err)
		return
	}

	domList := make([]DomainResult, 0)
	for _, dom := range domainList {
		var tenants []string
		for _, t := range dom.Tenants {
			tenants = append(tenants, t.String())
		}

		domRes := DomainResult{Domain: dom.Name.String(), Uuid: dom.Uuid.String(), Tenants: tenants}
		domList = append(domList, domRes)
	}

	result := DomainListResult{Domains: domList}
	c.JSON(http.StatusOK, result)
}

// Update handler doc
// @Tags Domain
// @Summary Update domain
// @Description Update domain info
// @Produce json
// @Param Tenant header string true "Tenant UUID to set access control"
// @Param domain_uuid path string true "Target domain's UUID"
// @Param domain body DomainUpdateRequest true "Request body parameter with json format"
// @Success 200 {object} DomainInfoResult
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/domains/{domain_uuid} [patch]
func (d *DomainController) Update(c Context) {
	requestTenant := c.GetHeader("Tenant")
	if requestTenant == "" {
		NewError(c,
			http.StatusBadRequest,
			errors.New("tenant uuid header is not specified"))
		return
	}

	requestTenantUuid, err := model.NewUuid(requestTenant)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		return
	}
	domainUuid := c.Param("domain_uuid")
	targetDomainUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Print(err)
		return
	}

	var request DomainUpdateRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Print(err)
		return
	}
	tenantList := request.Tenants
	if len(tenantList) == 0 {
		NewError(c, http.StatusBadRequest,
			errors.New("empty body parameter is given"))
	}

	var tenantUuidList []model.Uuid
	for _, t := range tenantList {
		tUuid, err := model.NewUuid(t)
		if err != nil {
			NewError(c, http.StatusBadRequest,
				errors.New("empty body parameter is given"))
		}
		tenantUuidList = append(tenantUuidList, tUuid)
	}

	domain, err := d.interactor.Update(targetDomainUuid, requestTenantUuid, tenantUuidList)
	if err != nil {
		switch e := err.(type) {
		case *model.InvalidParameterGiven, *model.DomainPermissionError:
			NewError(c, http.StatusBadRequest, err)
		case *model.DomainNotFoundError:
			NewError(c, http.StatusNotFound, err)
		default:
			NewError(c,
				http.StatusInternalServerError,
				NewUnAvailableHandlingError())
			log.Print(e)
		}
		log.Print(err)
		return
	}

	hosts := make([]HostResult, 0)
	for _, h := range domain.Hosts {
		host := HostResult{Name: h.Name, Address: h.Address, Uuid: h.Uuid.String()}
		hosts = append(hosts, host)
	}

	tenants := make([]string, 0)
	for _, t := range domain.Tenants {
		tenants = append(tenants, t.String())
	}

	var result DomainInfoResult
	result.Domain = domain.Name.String()
	result.Uuid = domain.Uuid.String()
	result.Tenants = tenants
	result.Hosts = hosts
	c.JSON(http.StatusOK, result)
}

// Get handler doc
// @Tags Domain
// @Summary Get domain
// @Description Get domain from coredns
// @Produce json
// @Param Tenant header string true "Tenant UUID to set access control"
// @Param domain_uuid path string true "Target domain's UUID"
// @Success 200 {object} DomainInfoResult
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/domains/{domain_uuid} [get]
func (d *DomainController) Get(c Context) {
	requestTenant := c.GetHeader("Tenant")
	if requestTenant == "" {
		NewError(c,
			http.StatusBadRequest,
			errors.New("tenant uuid header is not specified"))
		return
	}

	if requestTenant == "" {
		NewError(c,
			http.StatusBadRequest,
			errors.New("tenant uuid header is not specified"))
		return
	}
	requestTenantUuid, err := model.NewUuid(requestTenant)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		return
	}
	domainUuid := c.Param("domain_uuid")
	targetDomainUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Print(err)
		return
	}

	gotDomain, err := d.interactor.Get(targetDomainUuid, requestTenantUuid)
	if err != nil {
		switch e := err.(type) {
		case *model.InvalidParameterGiven, *model.DomainPermissionError:
			NewError(c, http.StatusBadRequest, err)
		case *model.DomainNotFoundError:
			NewError(c, http.StatusNotFound, err)
		default:
			NewError(c,
				http.StatusInternalServerError,
				NewUnAvailableHandlingError())
			log.Print(e)
		}
		log.Print(err)
		return
	}

	hosts := make([]HostResult, 0)
	for _, h := range gotDomain.Hosts {
		host := HostResult{Name: h.Name, Address: h.Address, Uuid: h.Uuid.String()}
		hosts = append(hosts, host)
	}

	tenants := make([]string, 0)
	for _, t := range gotDomain.Tenants {
		tenants = append(tenants, t.String())
	}

	var result DomainInfoResult
	result.Domain = gotDomain.Name.String()
	result.Uuid = gotDomain.Uuid.String()
	result.Tenants = tenants
	result.Hosts = hosts
	c.JSON(http.StatusOK, result)
}

// Delete handler doc
// @Tags Domain
// @Summary Delete domain
// @Description Delete new domain to coredns
// @Param Tenant header string true "Tenant UUID to set access control"
// @Param domain_uuid path string true "Target domain's UUID"
// @Success 204
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/domains/{domain_uuid} [delete]
func (d *DomainController) Delete(c Context) {
	requestTenant := c.GetHeader("Tenant")
	if requestTenant == "" {
		NewError(c,
			http.StatusBadRequest,
			errors.New("tenant uuid header is not specified"))
		return
	}

	requestTenantUuid, err := model.NewUuid(requestTenant)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		return
	}
	domainUuid := c.Param("domain_uuid")
	targetDomainUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Print(err)
		return
	}

	err = d.interactor.Delete(targetDomainUuid, requestTenantUuid)
	if err != nil {
		switch e := err.(type) {
		case *model.InvalidParameterGiven, *model.DomainPermissionError:
			NewError(c, http.StatusBadRequest, err)
		case *model.DomainNotFoundError:
			NewError(c, http.StatusNotFound, err)
		default:
			NewError(c,
				http.StatusInternalServerError,
				NewUnAvailableHandlingError())
			log.Print(e)
		}
		log.Print(err)
		return
	}

	c.Status(http.StatusNoContent)
}

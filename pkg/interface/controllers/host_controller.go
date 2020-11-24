package controllers

import (
	"errors"
	"log"
	"net/http"

	"coredns_api/internal/model"
	"coredns_api/internal/usecase"
)

type HostRequest struct {
	Name    string `json:"hostname"`
	Address string `json:"address"`
}

type HostController struct {
	interactor *usecase.HostInteractor
}

func NewHostController(itr *usecase.HostInteractor) *HostController {
	return &HostController{itr}
}

// Add handler doc
// @Tags Host
// @Summary Add new host
// @Description Add new host to domain
// @Accept json
// @Produce json
// @Param Tenant header string true "Tenant UUID to set access control"
// @Param domain_uuid path string true "Target domain's UUID"
// @Param host body HostRequest true "Request body parameter with json format"
// @Success 201 {object} DomainInfoResult
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/domains/{domain_uuid}/hosts [post]
func (d *HostController) Add(c Context) {
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
	targetDomain, err := d.interactor.GetDomain(targetDomainUuid, requestTenantUuid)
	if err != nil {
		switch e := err.(type) {
		case *model.DomainPermissionError:
			NewError(c, http.StatusBadRequest, err)
		default:
			NewError(c, http.StatusNotFound, err)
			log.Print(e)
		}
		log.Print(err)
		return
	}

	var requestedHost HostRequest
	err = c.ShouldBindJSON(&requestedHost)
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Print(err)
		return
	}

	name := requestedHost.Name
	address := requestedHost.Address
	newHost, err := model.NewOriginalHost(name, address, targetDomain.Name)
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

	gotDomain, err := d.interactor.Add(newHost, targetDomainUuid, requestTenantUuid)
	if err != nil {
		switch e := err.(type) {
		case *usecase.HostDuplicatedError:
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
		hr := HostResult{Name: h.Name, Address: h.Address, Uuid: h.Uuid.String()}
		hosts = append(hosts, hr)
	}
	hostRes := HostResult{Name: newHost.Name, Address: newHost.Address, Uuid: newHost.Uuid.String()}
	hosts = append(hosts, hostRes)

	var result DomainInfoResult
	result.Domain = gotDomain.Name.String()
	result.Uuid = gotDomain.Uuid.String()
	result.Hosts = hosts
	c.JSON(http.StatusCreated, result)
}

// List handler doc
// @Tags Host
// @Summary List hosts
// @Description List hosts from domain
// @Produce json
// @Param Tenant header string true "Tenant UUID to set access control"
// @Param domain_uuid path string true "Target domain's UUID"
// @Success 200 {object} DomainInfoResult
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/domains/{domain_uuid}/hosts [get]
func (d *HostController) List(c Context) {
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

	gotDomain, err := d.interactor.GetDomain(targetDomainUuid, requestTenantUuid)
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

	var result DomainInfoResult
	result.Domain = gotDomain.Name.String()
	result.Uuid = gotDomain.Uuid.String()
	result.Hosts = hosts
	c.JSON(http.StatusOK, result)
}

// Update handler doc
// @Tags Host
// @Summary Update host
// @Description Update host info
// @Accept json
// @Produce json
// @Param Tenant header string true "Tenant UUID to set access control"
// @Param domain_uuid path string true "Target domain's UUID"
// @Param host body HostRequest true "Request body parameter with json format"
// @Success 204 {object} DomainInfoResult
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/domains/{domain_uuid}/hosts [patch]
func (d *HostController) Update(c Context) {
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
	targetDomain, err := d.interactor.GetDomain(targetDomainUuid, requestTenantUuid)
	if err != nil {
		switch e := err.(type) {
		case *model.InvalidParameterGiven, *model.DomainPermissionError:
			NewError(c, http.StatusBadRequest, err)
		default:
			NewError(c, http.StatusBadRequest, err)
			log.Print(e)
		}
		log.Print(err)
		return
	}

	hostUuid := c.Param("host_uuid")
	targetHostUuid, err := model.NewUuid(hostUuid)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Print(err)
		return
	}

	host, err := d.interactor.Get(targetHostUuid, targetDomainUuid, requestTenantUuid)
	if err != nil {
		switch e := err.(type) {
		case *model.HostNotFoundError:
			NewError(c, http.StatusNotFound, err)
		case *model.DomainNotFoundError:
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

	var requestedHostInfo HostRequest
	err = c.ShouldBindJSON(&requestedHostInfo)
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Print(err)
		return
	}

	var name string
	if requestedHostInfo.Name != "" {
		name = requestedHostInfo.Name
	} else {
		name = host.Name
	}

	var address string
	if requestedHostInfo.Address != "" {
		address = requestedHostInfo.Address
	} else {
		address = host.Address
	}

	hostFqdn := model.GetFQDN(name, targetDomain.Name.String())
	updatedHost, err := model.NewHost(targetHostUuid, hostFqdn, address)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Print(err)
		return
	}

	err = d.interactor.Update(updatedHost, targetDomainUuid, requestTenantUuid)
	if err != nil {
		switch e := err.(type) {
		case *model.HostNotFoundError, *model.DomainNotFoundError:
			NewError(c, http.StatusNotFound, err)
		case *usecase.HostDuplicatedError, *model.DomainPermissionError:
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

	domain, err := d.interactor.GetDomain(targetDomainUuid, requestTenantUuid)
	if err != nil {
		switch e := err.(type) {
		case *model.DomainNotFoundError:
			NewError(c, http.StatusNotFound, err)
		case *model.DomainPermissionError:
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

	hosts := make([]HostResult, 0)
	for _, h := range domain.Hosts {
		hr := HostResult{Name: h.Name, Address: h.Address, Uuid: h.Uuid.String()}
		hosts = append(hosts, hr)
	}
	hostRes := HostResult{Name: updatedHost.Name, Address: updatedHost.Address, Uuid: updatedHost.Uuid.String()}
	hosts = append(hosts, hostRes)

	var result DomainInfoResult
	result.Domain = domain.Name.String()
	result.Uuid = domain.Uuid.String()
	result.Hosts = hosts
	c.JSON(http.StatusNoContent, result)
}

// Get handler doc
// @Tags Host
// @Summary Get host
// @Description Get host info
// @Produce json
// @Param Tenant header string true "Tenant UUID to set access control"
// @Param domain_uuid path string true "Target domain's UUID"
// @Param host_uuid path string true "Target host's UUID"
// @Success 200 {object} DomainInfoResult
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Router /v1/domains/{domain_uuid}/hosts/{host_uuid} [get]
func (d *HostController) Get(c Context) {
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

	hostUuid := c.Param("host_uuid")
	targetHostUuid, err := model.NewUuid(hostUuid)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Print(err)
		return
	}

	host, err := d.interactor.Get(targetHostUuid, targetDomainUuid, requestTenantUuid)
	if err != nil {
		switch e := err.(type) {
		case *model.HostNotFoundError, *model.DomainNotFoundError:
			NewError(c, http.StatusNotFound, err)
		case *model.DomainPermissionError:
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

	domain, err := d.interactor.GetDomain(targetDomainUuid, requestTenantUuid)
	if err != nil {
		switch e := err.(type) {
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
	hostRes := HostResult{Name: host.Name, Address: host.Address, Uuid: host.Uuid.String()}
	hosts = append(hosts, hostRes)

	var result DomainInfoResult
	result.Domain = domain.Name.String()
	result.Uuid = domain.Uuid.String()
	result.Hosts = hosts
	c.JSON(http.StatusOK, result)
}

// Delete handler doc
// @Tags Host
// @Summary Delete host
// @Description Delete host info
// @Param Tenant header string true "Tenant UUID to set access control"
// @Param domain_uuid path string true "Target domain's UUID"
// @Param host_uuid path string true "Target host's UUID"
// @Success 204 {object} DomainInfoResult
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Router /v1/domains/{domain_uuid}/hosts/{host_uuid} [delete]
func (d *HostController) Delete(c Context) {
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

	hostUuid := c.Param("host_uuid")
	targetHostUuid, err := model.NewUuid(hostUuid)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Print(err)
		return
	}

	host, err := d.interactor.Get(targetHostUuid, targetDomainUuid, requestTenantUuid)
	if err != nil {
		switch e := err.(type) {
		case *model.HostNotFoundError:
			NewError(c, http.StatusNotFound, err)
		case *model.DomainPermissionError:
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

	err = d.interactor.Delete(host, targetDomainUuid, requestTenantUuid)
	if err != nil {
		switch e := err.(type) {
		case *model.HostNotFoundError, *model.DomainNotFoundError:
			NewError(c, http.StatusNotFound, err)
		case *model.DomainPermissionError:
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

	c.Status(http.StatusNoContent)
}

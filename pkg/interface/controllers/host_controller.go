package controllers

import (
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
	Interactor *usecase.HostInteractor
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
// @Param domain_uuid path string true "domain_uuid"
// @Param hostname body string true "hostname"
// @Param address body string true "address"
// @Success 201 {object} AddResult
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/domains/{domain_uuid}/hosts [post]
func (d *HostController) Add(c Context) {
	domainUuid := c.Param("domain_uuid")
	targetDomainUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Fatal(err)
		return
	}

	var requestedHost HostRequest
	err = c.ShouldBindJSON(&requestedHost)
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Fatal(err)
		return
	}

	name := requestedHost.Name
	address := requestedHost.Address
	newHost, err := model.NewOriginalHost(name, address)
	if err != nil {
		switch e := err.(type) {
		case *model.InvalidParameterGiven:
			NewError(c, http.StatusBadRequest, err)
		default:
			NewError(c,
				http.StatusInternalServerError,
				NewUnAvailableHandlingError())
			log.Fatal(e)
		}
		log.Fatal(err)
		return
	}

	gotDomain, err := d.Interactor.Add(newHost, targetDomainUuid)
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Fatal(err)
		return
	}

	var hosts []HostResult
	hostRes := HostResult{Name: newHost.Name, Address: newHost.Address, Uuid: newHost.Uuid.String()}
	hosts = append(hosts, hostRes)

	var result AddResult
	result.Domain = gotDomain.Name.String()
	result.Uuid = gotDomain.Uuid.String()
	result.Hosts = hosts
	c.JSON(http.StatusCreated, result)
}

// Update handler doc
// @Tags Host
// @Summary Update host
// @Description Update host info
// @Accept json
// @Produce json
// @Param domain_uuid path string true "domain_uuid"
// @Param hostname body string false "hostname"
// @Param address body string false "address"
// @Success 204 {object} AddResult
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/domains/{domain_uuid}/hosts [patch]
func (d *HostController) Update(c Context) {
	domainUuid := c.Param("domain_uuid")
	targetDomainUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Fatal(err)
		return
	}

	hostUuid := c.Param("host_uuid")
	targetHostUuid, err := model.NewUuid(hostUuid)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Fatal(err)
		return
	}

	host, err := d.Interactor.Get(targetHostUuid, targetDomainUuid)
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Fatal(err)
		return
	}

	var requestedHostInfo HostRequest
	err = c.ShouldBindJSON(&requestedHostInfo)
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Fatal(err)
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

	updatedHost, err := model.NewHost(targetHostUuid, name, address)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Fatal(err)
		return
	}

	err = d.Interactor.Update(updatedHost, targetDomainUuid)
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Fatal(err)
		return
	}

	domain, _ := d.Interactor.GetDomain(targetDomainUuid)

	var hosts []HostResult
	hostRes := HostResult{Name: updatedHost.Name, Address: updatedHost.Address, Uuid: updatedHost.Uuid.String()}
	hosts = append(hosts, hostRes)

	var result AddResult
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
// @Param domain_uuid path string true "domain_uuid"
// @Param host_uuid path string true "host_uuid"
// @Success 200 {object} AddResult
// @Failure 400 {object} HTTPError
// @Router /v1/domains/{domain_uuid}/hosts/{host_uuid} [get]
func (d *HostController) Get(c Context) {
	domainUuid := c.Param("domain_uuid")
	targetDomainUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Fatal(err)
		return
	}

	hostUuid := c.Param("host_uuid")
	targetHostUuid, err := model.NewUuid(hostUuid)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Fatal(err)
		return
	}

	host, err := d.Interactor.Get(targetHostUuid, targetDomainUuid)
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Fatal(err)
		return
	}

	domain, err := d.Interactor.GetDomain(targetDomainUuid)
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Fatal(err)
		return
	}

	var hosts []HostResult
	hostRes := HostResult{Name: host.Name, Address: host.Address, Uuid: host.Uuid.String()}
	hosts = append(hosts, hostRes)

	var result AddResult
	result.Domain = domain.Name.String()
	result.Uuid = domain.Uuid.String()
	result.Hosts = hosts
	c.JSON(http.StatusOK, result)
}

// Delete handler doc
// @Tags Host
// @Summary Delete host
// @Description Delete host info
// @Param domain_uuid path string true "domain_uuid"
// @Param host_uuid path string true "host_uuid"
// @Success 204 {object} AddResult
// @Failure 400 {object} HTTPError
// @Router /v1/domains/{domain_uuid}/hosts/{host_uuid} [delete]
func (d *HostController) Delete(c Context) {
	domainUuid := c.Param("domain_uuid")
	targetDomainUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Fatal(err)
		return
	}

	hostUuid := c.Param("host_uuid")
	targetHostUuid, err := model.NewUuid(hostUuid)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Fatal(err)
		return
	}

	host, err := d.Interactor.Get(targetHostUuid, targetDomainUuid)
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Fatal(err)
		return
	}

	err = d.Interactor.Delete(host, targetDomainUuid)
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Fatal(err)
		return
	}

	c.Status(http.StatusNoContent)
}

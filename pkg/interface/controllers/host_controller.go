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
	dUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}

	var request HostRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		log.Fatal(err)
		mes := errors.New("unhandled server side error")
		NewError(c, http.StatusInternalServerError, mes)
		return
	}
	name := request.Name
	address := request.Address
	newHost, err := model.NewOriginalHost(name, address)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}

	gotDomain, err := d.Interactor.Add(newHost, dUuid)
	if err != nil {
		log.Fatal(err)
		mes := errors.New("unhandled server side error")
		NewError(c, http.StatusInternalServerError, mes)
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
	dUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}
	hostUuid := c.Param("host_uuid")
	hUuid, err := model.NewUuid(hostUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}

	var request HostRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		log.Fatal(err)
		mes := errors.New("unhandled server side error")
		NewError(c, http.StatusInternalServerError, mes)
		return
	}

	host, err := d.Interactor.Get(hUuid, dUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}

	var name string
	if request.Name != "" {
		name = request.Name
	} else {
		name = host.Name
	}

	var address string
	if request.Address != "" {
		address = request.Address
	} else {
		address = host.Address
	}

	newHost, err := model.NewHost(hUuid, name, address)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}

	err = d.Interactor.Update(newHost, dUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}

	domain, _ := d.Interactor.GetDomain(dUuid)

	var hosts []HostResult
	hostRes := HostResult{Name: newHost.Name, Address: newHost.Address, Uuid: newHost.Uuid.String()}
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
	dUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}
	hostUuid := c.Param("host_uuid")
	hUuid, err := model.NewUuid(hostUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}

	host, err := d.Interactor.Get(hUuid, dUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}

	domain, err := d.Interactor.GetDomain(dUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
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
	dUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}
	hostUuid := c.Param("host_uuid")
	hUuid, err := model.NewUuid(hostUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}

	host, err := d.Interactor.Get(hUuid, dUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}

	err = d.Interactor.Delete(host, dUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

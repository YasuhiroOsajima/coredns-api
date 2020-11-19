package controllers

import (
	"errors"
	"log"
	"net/http"

	"coredns_api/internal/model"
	"coredns_api/internal/usecase"
)

type DomainRequest struct {
	Name string `json:"domain"`
}

type AddResult struct {
	DomainResult
	Hosts []HostResult `json:"hosts"`
}

type DomainResult struct {
	Domain string `json:"domain"`
	Uuid   string `json:"uuid"`
}

type HostResult struct {
	Name    string `json:"hostname"`
	Address string `json:"address"`
	Uuid    string `json:"uuid"`
}

type DomainListResult struct {
	Domains []DomainResult `json:"domains"`
}

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
// @Param domain body string true "domain"
// @Success 201 {object} AddResult
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/domains [post]
func (d *DomainController) Add(c Context) {
	var request DomainRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		log.Fatal(err)
		mes := errors.New("unhandled server side error")
		NewError(c, http.StatusInternalServerError, mes)
		return
	}

	name := request.Name
	dom, err := model.NewOriginalDomain(name)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}

	err = d.interactor.Add(dom)
	if err != nil {
		log.Fatal(err)
		mes := errors.New("unhandled server side error")
		NewError(c, http.StatusInternalServerError, mes)
		return
	}

	var result AddResult
	result.Domain = dom.Name.String()
	result.Uuid = dom.Uuid.String()
	c.JSON(http.StatusCreated, result)
}

// List handler doc
// @Tags Domain
// @Summary List domains
// @Description List domains on coredns
// @Produce json
// @Success 200 {object} DomainListResult
// @Failure 500 {object} HTTPError
// @Router /v1/domains [get]
func (d *DomainController) List(c Context) {
	domainList, err := d.interactor.DbRepository.GetDomainsList()
	if err != nil {
		log.Fatal(err)
		mes := errors.New("unhandled server side error")
		NewError(c, http.StatusInternalServerError, mes)
		return
	}

	var domList []DomainResult
	for _, dom := range domainList {
		domRes := DomainResult{Domain: dom.Name.String(), Uuid: dom.Uuid.String()}
		domList = append(domList, domRes)
	}

	result := DomainListResult{Domains: domList}
	c.JSON(http.StatusOK, result)
}

// Get handler doc
// @Tags Domain
// @Summary Get new domain
// @Description Get new domain to coredns
// @Produce json
// @Param domain_uuid path string true "domain_uuid"
// @Success 200 {object} AddResult
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/domains/{domain_uuid} [get]
func (d *DomainController) Get(c Context) {
	domainUuid := c.Param("domain_uuid")
	dUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}

	gotDomain, err := d.interactor.Get(dUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusInternalServerError, err)
		return
	}

	var hosts []HostResult
	for _, h := range gotDomain.Hosts {
		host := HostResult{Name: h.Name, Address: h.Address, Uuid: h.Uuid.String()}
		hosts = append(hosts, host)
	}

	var result AddResult
	result.Domain = gotDomain.Name.String()
	result.Uuid = gotDomain.Uuid.String()
	result.Hosts = hosts
	c.JSON(http.StatusOK, result)
}

// Delete handler doc
// @Tags Domain
// @Summary Delete new domain
// @Description Delete new domain to coredns
// @Param domain_uuid path string true "domain_uuid"
// @Success 204
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/domains/{domain_uuid} [delete]
func (d *DomainController) Delete(c Context) {
	domainUuid := c.Param("domain_uuid")
	dUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusBadRequest, err)
		return
	}

	err = d.interactor.Delete(dUuid)
	if err != nil {
		log.Fatal(err)
		NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

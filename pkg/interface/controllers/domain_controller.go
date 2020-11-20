package controllers

import (
	"log"
	"net/http"

	"coredns_api/internal/model"
	"coredns_api/internal/usecase"
)

// Request
type DomainRequest struct {
	Name string `json:"domain"`
}

// Result
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
// @Param domain body string true "domain"
// @Success 201 {object} AddResult
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
	newDomain, err := model.NewOriginalDomain(name)
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

	var result AddResult
	result.Domain = newDomain.Name.String()
	result.Uuid = newDomain.Uuid.String()
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
	domainList, err := d.interactor.GetDomainsList()
	if err != nil {
		NewError(c,
			http.StatusInternalServerError,
			NewUnAvailableHandlingError())
		log.Print(err)
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
	targetDomainUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Print(err)
		return
	}

	gotDomain, err := d.interactor.Get(targetDomainUuid)
	if err != nil {
		switch e := err.(type) {
		case *model.InvalidParameterGiven:
			NewError(c, http.StatusBadRequest, err)
		case *usecase.RecordNotFoundError:
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
	targetDomainUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		log.Print(err)
		return
	}

	err = d.interactor.Delete(targetDomainUuid)
	if err != nil {
		switch e := err.(type) {
		case *model.InvalidParameterGiven:
			NewError(c, http.StatusBadRequest, err)
		case *usecase.RecordNotFoundError:
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

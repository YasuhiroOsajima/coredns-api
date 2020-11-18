package controllers

import (
	"errors"
	"log"

	"coredns_api/internal/model"
	"coredns_api/internal/usecase"
)

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

type DomainRequest struct {
	Name string `json:"domain"`
}

type DomainController struct {
	interactor *usecase.DomainInteractor
}

func NewDomainController(itr *usecase.DomainInteractor) *DomainController {
	return &DomainController{itr}
}

func (d *DomainController) Add(c Context) {
	var request DomainRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		log.Fatal(err)
		mes := errors.New("unhandled server side error")
		c.JSON(500, NewError(mes))
		return
	}

	name := request.Name
	dom, err := model.NewOriginalDomain(name)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}

	err = d.interactor.Add(dom)
	if err != nil {
		log.Fatal(err)
		mes := errors.New("unhandled server side error")
		c.JSON(500, NewError(mes))
		return
	}

	var result AddResult
	result.Domain = dom.Name
	result.Uuid = dom.Uuid.String()
	c.JSON(200, result)
}

func (d *DomainController) List(c Context) {
	domainList, err := d.interactor.DbRepository.GetDomainsList()
	if err != nil {
		log.Fatal(err)
		mes := errors.New("unhandled server side error")
		c.JSON(500, NewError(mes))
		return
	}

	var domList []DomainResult
	for _, dom := range domainList {
		domRes := DomainResult{Domain: dom.Name, Uuid: dom.Uuid.String()}
		domList = append(domList, domRes)
	}

	result := DomainListResult{Domains: domList}
	c.JSON(200, result)
}

// InstanceHandler is to get instances info
// @Summary Get instances info
// @Description A list of instances
// @Accept  plain
// @Produce  json
// @Success 200
// @Router /instances [get]
func (d *DomainController) Get(c Context) {
	domainUuid := c.Param("uuid")
	dUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}

	gotDomain, err := d.interactor.Get(dUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(500, NewError(err))
		return
	}

	var hosts []HostResult
	for _, h := range gotDomain.Hosts {
		host := HostResult{Name: h.Name, Address: h.Address, Uuid: h.Uuid.String()}
		hosts = append(hosts, host)
	}

	var result AddResult
	result.Domain = gotDomain.Name
	result.Uuid = gotDomain.Uuid.String()
	result.Hosts = hosts
	c.JSON(200, result)
}

func (d *DomainController) Delete(c Context) {
	domainUuid := c.Param("uuid")
	dUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}

	err = d.interactor.Delete(dUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(500, NewError(err))
		return
	}

	c.Status(204)
}

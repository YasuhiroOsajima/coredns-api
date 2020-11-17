package controllers

import (
	"errors"
	"log"

	"github.com/google/uuid"

	"coredns_api/internal/model"
	"coredns_api/internal/usecase"
)

type AddResult struct {
	Domain string       `json:"domain"`
	Uuid   string       `json:"uuid"`
	Hosts  []HostResult `json:"hosts"`
}

type HostResult struct {
	Name    string `json:"hostname"`
	Address string `json:"address"`
	Uuid    string `json:"uuid"`
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
	u, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
		mes := errors.New("unhandled server side error")
		c.JSON(500, NewError(mes))
		return
	}

	var request DomainRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		log.Fatal(err)
		mes := errors.New("unhandled server side error")
		c.JSON(500, NewError(mes))
		return
	}

	name := request.Name
	dmn, err := model.NewDomain(u.String(), name)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}

	err = d.interactor.Add(dmn)
	if err != nil {
		log.Fatal(err)
		mes := errors.New("Unhandled server side error")
		c.JSON(500, NewError(mes))
		return
	}

	result := AddResult{Domain: "", Uuid: ""}
	c.JSON(200, result)
}

func (d *DomainController) List(c Context) {

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
	name := c.Param("domain")
	domain, err := model.NewDefaultDomain(domainUuid, name)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}

	gotDomain, err := d.interactor.Get(domain)
	if err != nil {
		log.Fatal(err)
		c.JSON(500, NewError(err))
		return
	}

	var hosts []HostResult
	for _, h := range gotDomain.Hosts {
		host := HostResult{Name: h.Name, Address: h.Address, Uuid: h.Uuid}
		hosts = append(hosts, host)
	}

	result := AddResult{Domain: gotDomain.Name, Uuid: gotDomain.Uuid, Hosts: hosts}
	c.JSON(200, result)
}

func (d *DomainController) Delete(c Context) {
	domainUuid := c.Param("uuid")
	name := c.Param("domain")
	domain, err := model.NewDefaultDomain(domainUuid, name)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}

	err = d.interactor.Delete(domain)
	if err != nil {
		log.Fatal(err)
		c.JSON(500, NewError(err))
		return
	}

	c.Status(204)
}

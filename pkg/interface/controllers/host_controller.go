package controllers

import (
	"coredns_api/internal/model"
	"coredns_api/internal/usecase"
	"errors"
	"log"
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

// ImageHandler is to get images info
// @Summary Get images info
// @Description A list of images
// @Accept  plain
// @Produce  json
// @Success 200
// @Router /images [get]
func (d *HostController) Add(c Context) {
	domainUuid := c.Param("domain_uuid")
	dUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}

	var request HostRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		log.Fatal(err)
		mes := errors.New("unhandled server side error")
		c.JSON(500, NewError(mes))
		return
	}
	name := request.Name
	address := request.Address
	newHost, err := model.NewOriginalHost(name, address)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}

	gotDomain, err := d.Interactor.Add(newHost, dUuid)
	if err != nil {
		log.Fatal(err)
		mes := errors.New("unhandled server side error")
		c.JSON(500, NewError(mes))
		return
	}

	var hosts []HostResult
	hostRes := HostResult{Name: newHost.Name, Address: newHost.Address, Uuid: newHost.Uuid.String()}
	hosts = append(hosts, hostRes)

	var result AddResult
	result.Domain = gotDomain.Name.String()
	result.Uuid = gotDomain.Uuid.String()
	result.Hosts = hosts
	c.JSON(200, result)
}

func (d *HostController) Update(c Context) {
	domainUuid := c.Param("domain_uuid")
	dUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}
	hostUuid := c.Param("uuid")
	hUuid, err := model.NewUuid(hostUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}

	var request HostRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		log.Fatal(err)
		mes := errors.New("unhandled server side error")
		c.JSON(500, NewError(mes))
		return
	}

	host, err := d.Interactor.Get(hUuid, dUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
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
		c.JSON(400, NewError(err))
		return
	}

	err = d.Interactor.Update(newHost, dUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
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
	c.JSON(200, result)
}

func (d *HostController) Get(c Context) {
	domainUuid := c.Param("domain_uuid")
	dUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}
	hostUuid := c.Param("uuid")
	hUuid, err := model.NewUuid(hostUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}

	host, err := d.Interactor.Get(hUuid, dUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}

	domain, err := d.Interactor.GetDomain(dUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}

	var hosts []HostResult
	hostRes := HostResult{Name: host.Name, Address: host.Address, Uuid: host.Uuid.String()}
	hosts = append(hosts, hostRes)

	var result AddResult
	result.Domain = domain.Name.String()
	result.Uuid = domain.Uuid.String()
	result.Hosts = hosts
	c.JSON(200, result)
}

func (d *HostController) Delete(c Context) {
	domainUuid := c.Param("domain_uuid")
	dUuid, err := model.NewUuid(domainUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}
	hostUuid := c.Param("uuid")
	hUuid, err := model.NewUuid(hostUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}

	host, err := d.Interactor.Get(hUuid, dUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}

	err = d.Interactor.Delete(host, dUuid)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, NewError(err))
		return
	}

	c.Status(204)
}

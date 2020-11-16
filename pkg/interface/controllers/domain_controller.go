package controllers

import (
	"coredns_api/internal/usecase"
)

type AddResult struct {
	Domain string `json:"domain"`
	Uuid   string `json:"uuid"`
}

type DomainController struct {
	interactor usecase.DomainInteractor
}

func NewDomainController(itr usecase.DomainInteractor) *DomainController {
	return &DomainController{itr}
}

// InstanceHandler is to get instances info
// @Summary Get instances info
// @Description A list of instances
// @Accept  plain
// @Produce  json
// @Success 200
// @Router /instances [get]
// func (d *DomainController) Add(c Context) {
// 	u, err := uuid.NewRandom()
// 	if err != nil {
// 		log.Fatal(err)
// 		mes := errors.New("Unhandled server side error")
// 		c.JSON(500, NewError(mes))
// 		return
// 	}
//
// 	name := c.Param("domain")
// 	dmn, err := model.NewDomain(u.String(), name)
// 	if err != nil {
// 		log.Fatal(err)
// 		c.JSON(400, NewError(err))
// 		return
// 	}
//
// 	err = d.interactor.Add(dmn)
// 	if err != nil {
// 		log.Fatal(err)
// 		mes := errors.New("Unhandled server side error")
// 		c.JSON(500, NewError(mes))
// 		return
// 	}
//
// 	result := AddResult{domain: "", uuid: ""}
// 	c.JSON(200, result)
// }

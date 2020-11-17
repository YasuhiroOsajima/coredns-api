package controllers

import (
	"coredns_api/internal/usecase"
)

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
	result := AddResult{Domain: "", Uuid: ""}
	c.JSON(200, result)
}

func (d *HostController) Update(c Context) {
	result := AddResult{Domain: "", Uuid: ""}
	c.JSON(200, result)
}

func (d *HostController) Get(c Context) {
	result := AddResult{Domain: "", Uuid: ""}
	c.JSON(200, result)
}

func (d *HostController) Delete(c Context) {
	result := AddResult{Domain: "", Uuid: ""}
	c.JSON(200, result)
}

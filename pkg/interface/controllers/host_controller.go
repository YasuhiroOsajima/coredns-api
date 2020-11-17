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

// func HostsHndler(c Context) {
// 	db := repository.NewImageRepository()
// 	imagelist, err := db.FindAll()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	var images []Image
// 	for _, i := range imagelist {
// 		var img Image
// 		img.Uuid = i.Uuid
// 		img.Name = i.Name
// 		img.Owner = i.Owner
// 		images = append(images, img)
// 	}
//
// 	c.JSON(200, images)
// }

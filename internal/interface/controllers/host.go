package controllers

import (
	"log"

	"github.com/gin-gonic/gin"

	"swaggo_sample/internal/repository"
)

type Image struct {
	Uuid  string `json:"uuid"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

// ImageHandler is to get images info
// @Summary Get images info
// @Description A list of images
// @Accept  plain
// @Produce  json
// @Success 200
// @Router /images [get]
func ImageHandler(c *gin.Context) {
	db := repository.NewImageRepository()
	imagelist, err := db.FindAll()
	if err != nil {
		log.Fatal(err)
	}

	var images []Image
	for _, i := range imagelist {
		var img Image
		img.Uuid = i.Uuid
		img.Name = i.Name
		img.Owner = i.Owner
		images = append(images, img)
	}

	c.JSON(200, images)
}

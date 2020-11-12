package infrastructure

import (
	"os"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "coredns_api/docs"
	c "coredns_api/internal/interface/controllers"
)

// @title gin-swagger sample
// @version 0.1
func Router() {
	var Server = os.Getenv("SERVER")
	var Port = os.Getenv("PORT")

	var Router *gin.Engine
	Router = gin.Default()

	Router.GET("/instances", c.InstanceHandler)
	Router.GET("/images", c.ImageHandler)

	url := ginSwagger.URL("http://" + Server + ":" + Port + "/swagger/doc.json")
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	Router.Run(":" + Port)
}

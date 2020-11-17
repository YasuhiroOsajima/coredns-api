package infrastructure

import (
	"os"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	// _ "coredns_api/docs"
)

// @title gin-swagger sample
// @version 0.1
func Router() {
	dcntr := InitializeDomainController()
	hcntr := InitializeHostController()

	var Server = os.Getenv("SERVER")
	var Port = os.Getenv("PORT")

	var Router *gin.Engine
	Router = gin.Default()

	Router.POST("/domains", func(c *gin.Context) { dcntr.Add(c) })
	Router.GET("/domains/:uuid", func(c *gin.Context) { dcntr.Get(c) })

	Router.GET("/hosts", func(c *gin.Context) { hcntr.Add(c) })

	url := ginSwagger.URL("http://" + Server + ":" + Port + "/swagger/doc.json")
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	Router.Run(":" + Port)
}

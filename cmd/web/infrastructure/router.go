package infrastructure

import (
	"os"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "coredns_api/docs"
)

func Router() {
	dcntr := InitializeDomainController()
	hcntr := InitializeHostController()

	var Server = os.Getenv("SERVER")
	var Port = os.Getenv("PORT")

	var Router *gin.Engine
	Router = gin.Default()

	Router.POST("/v1/domains", func(c *gin.Context) { dcntr.Add(c) })
	Router.GET("/v1/domains", func(c *gin.Context) { dcntr.List(c) })
	Router.GET("/v1/domains/:domain_uuid", func(c *gin.Context) { dcntr.Get(c) })
	Router.PATCH("/v1/domains/:domain_uuid", func(c *gin.Context) { dcntr.Update(c) })
	Router.DELETE("/v1/domains/:domain_uuid", func(c *gin.Context) { dcntr.Delete(c) })

	Router.POST("/v1/domains/:domain_uuid/hosts", func(c *gin.Context) { hcntr.Add(c) })
	Router.GET("/v1/domains/:domain_uuid/hosts", func(c *gin.Context) { hcntr.List(c) })
	Router.PATCH("/v1/domains/:domain_uuid/hosts/:host_uuid", func(c *gin.Context) { hcntr.Update(c) })
	Router.GET("/v1/domains/:domain_uuid/hosts/:host_uuid", func(c *gin.Context) { hcntr.Get(c) })
	Router.DELETE("/v1/domains/:domain_uuid/hosts/:host_uuid", func(c *gin.Context) { hcntr.Delete(c) })

	url := ginSwagger.URL("http://" + Server + ":" + Port + "/swagger/doc.json")
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	Router.Run(":" + Port)
}

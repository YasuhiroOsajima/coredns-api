package infrastructure

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"os"

	_ "coredns_api/docs"
)

// Todo. メモリ上にドメイン情報を置いておく
// Todo. Goルーチンで定期的に書き出し

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
	Router.GET("/domains", func(c *gin.Context) { dcntr.List(c) })
	Router.GET("/domains/:uuid", func(c *gin.Context) { dcntr.Get(c) })
	Router.DELETE("/domains/:uuid", func(c *gin.Context) { dcntr.Delete(c) })

	Router.POST("/domains/:domain_uuid/hosts", func(c *gin.Context) { hcntr.Add(c) })
	Router.PATCH("/domains/:domain_uuid/hosts/:uuid", func(c *gin.Context) { hcntr.Update(c) })
	Router.GET("/domains/:domain_uuid/hosts/:uuid", func(c *gin.Context) { hcntr.Get(c) })
	Router.DELETE("/domains/:domain_uuid/hosts/:uuid", func(c *gin.Context) { hcntr.Delete(c) })

	url := ginSwagger.URL("http://" + Server + ":" + Port + "/swagger/doc.json")
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	Router.Run(":" + Port)
}

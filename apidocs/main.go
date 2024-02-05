package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaoyan648/learn/apidocs/controller"
	"log"
)
import ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
import swaggerFiles "github.com/swaggo/files"     // swagger embed files

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	r := gin.Default()

	ctr := &controller.Controller{}

	v1 := r.Group("/api/v1")
	{
		accounts := v1.Group("/accounts")
		{
			accounts.GET(":id", ctr.ShowAccount)
			accounts.GET("", ctr.ListAccounts)
		}
		//...
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

//...

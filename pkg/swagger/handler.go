package swagger

import (
	"strings"

	"ienergy-template-go/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type swagger struct {
}

func NewSwagger() *swagger {
	return &swagger{}
}

func (m *swagger) Register(gGroup gin.IRouter) {
	g := gGroup.Group("")
	{
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
		g.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

func (m *swagger) SwaggerHandler(isProduction bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isProduction {
			return
		}
		docs.SwaggerInfo.Title = "ienergy-template-go swagger"
		docs.SwaggerInfo.Description = "API Documentation for ienergy-template-go"
		docs.SwaggerInfo.Host = strings.ToLower(c.Request.Host)
		docs.SwaggerInfo.BasePath = "/api/v1"
		c.Next()
	}
}

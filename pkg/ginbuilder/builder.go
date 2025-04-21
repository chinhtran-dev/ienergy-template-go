package ginbuilder

import (
	"github.com/gin-gonic/gin"
)

type builder struct {
	middlewares []gin.HandlerFunc
}

// BaseBuilder initializes a builder with default middleware like gin.Recovery() to handle panic
func BaseBuilder() *builder {
	return &builder{
		middlewares: []gin.HandlerFunc{
			gin.Recovery(),
		},
	}
}

// Default initializes an empty builder without any middlewares
func Default() *builder {
	return &builder{}
}

// Build
func (b *builder) Build() *gin.Engine {
	e := defaultGinEngine()
	e.Use(b.middlewares...)
	return e
}

func defaultGinEngine() *gin.Engine {
	e := gin.New()
	return e
}

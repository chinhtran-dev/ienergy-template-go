package util

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
*
This feature should be use for gin middleware that binding multiple times without any damages.
This will restore body into the request after binding
so you should not use this method if you are
enough to call binding at once.
*/
func ShouldBind(c *gin.Context, obj interface{}) error {
	if c.Request.Method == http.MethodGet {
		return c.Bind(obj)
	}
	if c.Request.Body == nil {
		return nil
	}
	var bodyBytes []byte
	bodyBytes, _ = io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := c.Bind(obj); err != nil {
		return err
	}

	// write back to request body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return nil
}

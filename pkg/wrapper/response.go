package wrapper

import (
	"encoding/json"
	"ienergy-template-go/pkg/errormap"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int         `json:"status_code"`
	Code       int         `json:"code"`
	Data       interface{} `json:"data,omitempty"`
	Message    *string     `json:"message,omitempty"`
	TrackID    string      `json:"trace_id"`
}

func (r *Response) String() string {
	data, _ := json.Marshal(r) // nolint:errchkjson
	return string(data)
}

func NewResponse(c *gin.Context, statusCode int, data interface{}, message *string, err error) *Response {
	if message == nil {
		newmessage := errormap.ErrorMapMsg[statusCode]
		message = &newmessage
	}
	// overwrite error if exists
	if err != nil {
		messageError := err.Error()
		message = &messageError
	}
	return &Response{
		Data:       data,
		Message:    message,
		Code:       errormap.ErrorMapCode[statusCode],
		StatusCode: statusCode,
	}
}

func NewSuccessResponse(c *gin.Context, data interface{}) *Response {
	return NewResponse(c, http.StatusOK, data, nil, nil)
}

func NewNotFoundResponse(c *gin.Context, data interface{}, err error) *Response {
	return NewResponse(c, http.StatusNotFound, data, nil, err)
}

func NewUnauthorizationResponse(c *gin.Context, data interface{}, err error) *Response {
	return NewResponse(c, http.StatusUnauthorized, data, nil, err)
}

func NewBadRequestResponse(c *gin.Context, data interface{}, err error) *Response {
	return NewResponse(c, http.StatusBadRequest, data, nil, err)
}

func NewInternalServerErrorResponse(c *gin.Context, data interface{}, err error) *Response {
	return NewResponse(c, http.StatusInternalServerError, data, nil, err)
}

func JSON401(c *gin.Context, data interface{}, err error) {
	c.JSON(http.StatusUnauthorized, NewUnauthorizationResponse(c, data, err))
}

func JSON404(c *gin.Context, data interface{}, err error) {
	c.JSON(http.StatusNotFound, NewNotFoundResponse(c, data, err))
}

func JSON400(c *gin.Context, data interface{}, err error) {
	c.JSON(http.StatusBadRequest, NewBadRequestResponse(c, data, err))
}

func JSON500(c *gin.Context, data interface{}, err error) {
	c.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse(c, data, err))
}

func JSONOk(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, NewSuccessResponse(c, data))
}

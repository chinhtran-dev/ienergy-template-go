package handler

import (
	"ienergy-template-go/internal/model/request"
	"ienergy-template-go/internal/service"
	"ienergy-template-go/pkg/wrapper"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(
	authService service.AuthService,
) AuthHandler {
	return AuthHandler{
		authService: authService,
	}
}

// User godoc
// @Summary Get my record page for non login usersas
// @Description Get my record page for non login usersas
// @Tags auth
// @Accept json
// @Produce json
// @Param model body request.UserRegisterRequest true "model"
// @Success 200 {object} wrapper.Response
// @Failure 400 {object} wrapper.Response
// @Failure 500 {object} wrapper.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.UserRegisterRequest
		if err := c.BindJSON(&req); err != nil {
			c.Error(err)
			return
		}
		err := req.Validate()
		if err != nil {
			c.Error(err)
			return
		}
		resp, err := h.authService.Register(c, req)
		if err != nil {
			c.Error(err)
			return
		}
		wrapper.JSONOk(c, resp)
	}
}

// User godoc
// @Summary API for get token from user name email and password
// @Description API for get token from user name email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param model body request.UserLoginRequest true "model"
// @Success 200 {object} wrapper.Response{data=string} "success"
// @Failure 400 {object} wrapper.Response
// @Failure 500 {object} wrapper.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.UserLoginRequest
		if err := c.BindJSON(&req); err != nil {
			c.Error(err)
			return
		}
		err := req.Validate()
		if err != nil {
			c.Error(err)
			return
		}
		resp, err := h.authService.Login(c, req)

		if err != nil {
			c.Error(err)
			return
		}

		if len(resp.Token) == 0 {
			c.Error(err)
			return
		}
		wrapper.JSONOk(c, resp)
	}
}

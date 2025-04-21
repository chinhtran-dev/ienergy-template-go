package handler

import (
	"ienergy-template-go/internal/service"
	"ienergy-template-go/pkg/wrapper"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{
		userService: userService,
	}
}

// User godoc
// @Summary API for get token from user name email and password
// @Description API for get token from user name email and password
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} wrapper.Response{data=response.UserInfoResponse} "success"
// @Failure 400 {object} wrapper.Response
// @Failure 500 {object} wrapper.Response
// @Router /user/info [get]
func (h *UserHandler) Info() gin.HandlerFunc {
	return func(c *gin.Context) {
		info, err := h.userService.GetUserInfo(c)
		if err != nil {
			wrapper.JSON401(c, nil, err)
			return
		}

		wrapper.JSONOk(c, info)
	}
}

package handler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	us domain.UserService
}

func NewUserHandler(rg *gin.RouterGroup, userService *domain.UserService) {
	uh := &UserHandler{
		us: *userService,
	}

	us := rg.Group("users")

	us.GET("/me", uh.Me)
}

func (u *UserHandler) Me(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Hello me!")
}

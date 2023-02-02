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

func (u *UserHandler) Create(ctx *gin.Context) {

	var user domain.User
	if err := ctx.BindJSON(&user); err != nil {
		newError := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newError), gin.H{"error": newError})
		return
	}

	if err := u.us.Store(ctx, &user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}

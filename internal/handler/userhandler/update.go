package userhandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type updateUser struct {
	Password     string `json:"password" binding:"required"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Permission   string `json:"permission" binding:"required"`
	ProfileColor string `json:"profile_color" binding:"required"`
}

func (uh *UserHandler) UpdateByID(ctx *gin.Context) {
	var uUser updateUser
	if err := ctx.BindJSON(&uUser); err != nil {
		newError := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newError), gin.H{"error": newError})
		return
	}

	user := domain.User{
		Password:     uUser.Password,
		FirstName:    uUser.FirstName,
		LastName:     uUser.LastName,
		Permission:   uUser.Permission,
		ProfileColor: uUser.ProfileColor,
	}

	context := ctx.Request.Context()
	if err := uh.userService.Update(context, &user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

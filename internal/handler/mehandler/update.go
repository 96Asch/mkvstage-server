package mehandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type updateUser struct {
	Password     string `json:"password" binding:"required"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	ProfileColor string `json:"profile_color" binding:"required"`
}

func (uh *meHandler) Update(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		err := domain.NewInternalErr()
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	var uUser updateUser
	if err := ctx.BindJSON(&uUser); err != nil {
		newError := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newError), gin.H{"error": newError})
		return
	}

	user := domain.User{
		ID:           val.(*domain.User).ID,
		Password:     uUser.Password,
		FirstName:    uUser.FirstName,
		LastName:     uUser.LastName,
		Permission:   val.(*domain.User).Permission,
		ProfileColor: uUser.ProfileColor,
	}

	context := ctx.Request.Context()
	if err := uh.userService.Update(context, &user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

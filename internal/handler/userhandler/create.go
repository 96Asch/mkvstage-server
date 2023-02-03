package userhandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

func (u *UserHandler) Create(ctx *gin.Context) {

	var user domain.User
	if err := ctx.BindJSON(&user); err != nil {
		newError := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newError), gin.H{"error": newError})
		return
	}

	if err := u.userService.Store(ctx, &user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}

package userhandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

func (u *userHandler) Me(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		err := domain.NewInternalErr()
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	context := ctx.Request.Context()
	id := val.(*domain.User).ID

	user, err := u.userService.FetchByID(context, id)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

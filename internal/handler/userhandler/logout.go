package userhandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

func (us userHandler) Logout(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": newErr})
		return
	}

	user := val.(*domain.User)

	context := ctx.Request.Context()
	err := us.tokenService.Logout(context, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.Status(http.StatusAccepted)
}

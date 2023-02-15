package mehandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type logoutReq struct {
	Refresh string `json:"refresh" binding:"required"`
}

func (us meHandler) Logout(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": newErr})
		return
	}

	var logout logoutReq
	if err := ctx.BindJSON(&logout); err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})
		return
	}

	user := val.(*domain.User)

	context := ctx.Request.Context()
	err := us.tokenService.RemoveRefresh(context, user.ID, logout.Refresh)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.Status(http.StatusAccepted)
}

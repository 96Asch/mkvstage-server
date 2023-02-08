package userhandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

func (us *UserHandler) GetAll(ctx *gin.Context) {
	context := ctx.Request.Context()
	users, err := us.userService.FetchAll(context)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

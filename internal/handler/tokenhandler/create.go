package tokenhandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type refreshRequest struct {
	Refresh string `json:"refresh" binding:"required"`
}

func (th tokenHandler) CreateAccess(ctx *gin.Context) {
	var req refreshRequest
	if err := ctx.BindJSON(&req); err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	context := ctx.Request.Context()

	access, err := th.tokenService.CreateAccess(context, req.Refresh)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": access})
}

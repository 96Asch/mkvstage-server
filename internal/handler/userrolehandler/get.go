package userrolehandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

func (urh userRoleHandler) Me(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})
		return
	}

	user := val.(*domain.User)

	context := ctx.Request.Context()
	userroles, err := urh.urs.FetchByUser(context, user)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user_roles": userroles})
}

func (urh userRoleHandler) GetAll(ctx *gin.Context) {
	context := ctx.Request.Context()
	userroles, err := urh.urs.FetchAll(context)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user_roles": userroles})
}

package mehandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

func (mh meHandler) Me(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		err := domain.NewInternalErr()
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	tokenUser, ok := val.(*domain.User)
	if !ok {
		err := domain.NewInternalErr()
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	context := ctx.Request.Context()

	user, err := mh.userService.FetchByID(context, tokenUser.ID)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

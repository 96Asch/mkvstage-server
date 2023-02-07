package userhandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type deleteID struct {
	ID int64 `json:"id"`
}

func (uh *UserHandler) Delete(ctx *gin.Context) {

	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})
		return
	}

	var dID deleteID
	if err := ctx.BindJSON(&dID); err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})
		return
	}

	user := val.(*domain.User)
	context := ctx.Request.Context()
	if err := uh.userService.Remove(context, user, dID.ID); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.Status(http.StatusAccepted)
}

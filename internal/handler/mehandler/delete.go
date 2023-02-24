package mehandler

import (
	"log"
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type deleteID struct {
	ID int64 `json:"id" binding:"required"`
}

func (mh meHandler) Delete(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})
		log.Println(newErr)

		return
	}

	var dID deleteID
	if err := ctx.BindJSON(&dID); err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})
		log.Println(newErr)

		return
	}

	user := val.(*domain.User)
	context := ctx.Request.Context()

	id, err := mh.userService.Remove(context, user, dID.ID)
	if err != nil {
		log.Println(err)
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	if err := mh.tokenService.RemoveAllRefresh(context, id); err != nil {
		log.Println(err)
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	ctx.Status(http.StatusAccepted)
}

package rolehandler

import (
	"net/http"
	"strconv"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

func (rh roleHandler) DeleteByID(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	idField := ctx.Params.ByName("id")

	roleID, err := strconv.Atoi(idField)
	if err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	context := ctx.Request.Context()

	user, ok := val.(*domain.User)
	if !ok {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	err = rh.rs.Remove(context, int64(roleID), user)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	ctx.Status(http.StatusAccepted)
}
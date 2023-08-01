package setlistrolehandler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

func (srh setlistRoleHandler) Delete(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(newErr.Status(), gin.H{"error": newErr.Error()})

		return
	}

	user, ok := val.(*domain.User)
	if !ok {
		newErr := domain.NewInternalErr()
		ctx.JSON(newErr.Status(), gin.H{"error": newErr.Error()})

		return
	}

	setlistIDsString := ctx.Query("ids")

	if setlistIDsString == "" {
		ctx.Status(http.StatusOK)

		return
	}

	splitSetlistIDs := strings.Split(setlistIDsString, ",")
	setlistIDs := make([]int64, len(splitSetlistIDs))

	for idx, id := range splitSetlistIDs {
		val, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Could not parse %s to a number", id)})

			return
		}

		setlistIDs[idx] = val
	}

	context := ctx.Request.Context()
	err := srh.slrs.Remove(context, setlistIDs, user)

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

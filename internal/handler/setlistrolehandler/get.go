package setlistrolehandler

import (
	"net/http"
	"strconv"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

func (srh setlistRoleHandler) GetAll(ctx *gin.Context) {
	querySetlistID := ctx.Query("setlist")

	var querySetlists *[]domain.Setlist

	if len(querySetlistID) != 0 {
		convSetlistID, err := strconv.Atoi(querySetlistID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		querySetlists = &[]domain.Setlist{{ID: int64(convSetlistID)}}
	}

	context := ctx.Request.Context()
	setlistRoles, err := srh.slrs.Fetch(context, querySetlists)

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"setlistroles": setlistRoles})
}

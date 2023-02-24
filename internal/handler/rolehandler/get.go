package rolehandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

func (rh roleHandler) GetAll(ctx *gin.Context) {
	context := ctx.Request.Context()

	roles, err := rh.rs.FetchAll(context)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"roles": roles})
}

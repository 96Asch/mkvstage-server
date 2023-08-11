package songhandler

import (
	"net/http"
	"strconv"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

func (sh songHandler) GetByID(ctx *gin.Context) {
	idField := ctx.Params.ByName("id")

	songID, err := strconv.Atoi(idField)
	if err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	context := ctx.Request.Context()

	song, err := sh.ss.FetchByID(context, int64(songID))
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"song": song})
}

func (sh songHandler) GetAll(ctx *gin.Context) {
	context := ctx.Request.Context()

	songs, err := sh.ss.FetchAll(context)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"songs": songs})
}

package setlisthandler

import (
	"net/http"
	"strconv"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/util"
	"github.com/gin-gonic/gin"
)

func (slh setlistHandler) GetAll(ctx *gin.Context) {
	fromTime, fromErr := util.StringToTime(ctx.Query("from"))
	toTime, toErr := util.StringToTime(ctx.Query("to"))

	if fromErr != nil {
		ctx.JSON(domain.Status(fromErr), gin.H{"error": fromErr.Error()})

		return
	}

	if toErr != nil {
		ctx.JSON(domain.Status(toErr), gin.H{"error": toErr.Error()})

		return
	}

	context := ctx.Request.Context()
	retrievedSetlists, err := slh.sls.Fetch(context, fromTime, toTime)

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"setlist": retrievedSetlists,
	})
}

func (slh setlistHandler) GetByID(ctx *gin.Context) {
	idField := ctx.Params.ByName("id")

	setlistID, err := strconv.Atoi(idField)
	if err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	context := ctx.Request.Context()
	setlist, err := slh.sls.FetchByID(context, int64(setlistID))

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	setlistEntries, err := slh.sles.FetchBySetlist(context, setlist)

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"setlist": setlist,
		"entries": setlistEntries,
	})
}

package setlisthandler

import (
	"net/http"

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

	retrievedSetlistEntries, err := slh.sles.FetchBySetlist(context, retrievedSetlists)

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"setlist": retrievedSetlists,
		"entries": retrievedSetlistEntries,
	})
}

func (slh setlistHandler) GetByID(ctx *gin.Context) {
	fields, err := util.BindNamedParams(ctx, "id")

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	context := ctx.Request.Context()
	setlist, err := slh.sls.FetchByID(context, fields["id"])

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	setlistEntries, err := slh.sles.FetchBySetlist(context, &[]domain.Setlist{*setlist})

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"setlist": setlist,
		"entries": setlistEntries,
	})
}

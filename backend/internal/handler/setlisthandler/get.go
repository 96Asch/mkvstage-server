package setlisthandler

import (
	"log"
	"net/http"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/96Asch/mkvstage-server/backend/internal/util"
	"github.com/gin-gonic/gin"
)

type setlistResponse struct {
	domain.Setlist
	Entries []domain.SetlistEntry `json:"entries"`
}

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

	response := make([]setlistResponse, len(*retrievedSetlists))
	sortedEntries := sortBySetlist(retrievedSetlistEntries)

	for idx, setlist := range *retrievedSetlists {
		log.Printf("sid: %d", setlist.ID)
		response[idx] = setlistResponse{
			setlist,
			sortedEntries[setlist.ID],
		}
		log.Println(response[idx])
	}

	ctx.JSON(http.StatusOK, gin.H{
		"setlists": response,
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

	response := setlistResponse{
		*setlist,
		*setlistEntries,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"setlist": response,
	})
}

func sortBySetlist(entries *[]domain.SetlistEntry) map[int64][]domain.SetlistEntry {
	sortedBySetlist := make(map[int64][]domain.SetlistEntry)

	for _, entry := range *entries {
		val := sortedBySetlist[entry.SetlistID]
		sortedBySetlist[entry.SetlistID] = append(val, entry)
	}

	return sortedBySetlist
}

package setlisthandler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type setlistUpdateReq struct {
	Name           string                `json:"name" binding:"required"`
	CreatorID      int64                 `json:"creator_id" binding:"required"`
	Deadline       time.Time             `json:"deadline" binding:"required"`
	CreatedEntries []domain.SetlistEntry `json:"created_entries" binding:"required"`
	UpdatedEntries []domain.SetlistEntry `json:"updated_entries" binding:"required"`
	DeletedEntries []domain.SetlistEntry `json:"deleted_entries" binding:"required"`
}

func (slh setlistHandler) UpdateByID(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	idField := ctx.Params.ByName("id")

	setlistID, err := strconv.Atoi(idField)
	if err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	var slReq setlistUpdateReq
	if err := ctx.BindJSON(&slReq); err != nil {
		newErr := domain.NewBadRequestErr(err.Error())
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	user, ok := val.(*domain.User)
	if !ok {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	setlist := &domain.Setlist{
		ID:        int64(setlistID),
		Name:      slReq.Name,
		CreatorID: slReq.CreatorID,
		Deadline:  slReq.Deadline.Local(),
	}

	createdEntries := make([]domain.SetlistEntry, len(slReq.CreatedEntries))
	updatedEntries := make([]domain.SetlistEntry, len(slReq.UpdatedEntries))
	deletedEntriesIds := make([]int64, len(slReq.DeletedEntries))

	copy(createdEntries, slReq.CreatedEntries)
	copy(updatedEntries, slReq.UpdatedEntries)

	for idx, entry := range slReq.DeletedEntries {
		deletedEntriesIds[idx] = entry.ID
	}

	context := ctx.Request.Context()
	updatedSetlist, err := slh.sls.Update(context, setlist, user)

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	if err := slh.sles.StoreBatch(context, &createdEntries, user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	if err := slh.sles.UpdateBatch(context, &updatedEntries, user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	if err := slh.sles.RemoveBatch(context, setlist, deletedEntriesIds, user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"setlist": updatedSetlist,
		"entries": append(createdEntries, updatedEntries...),
	})
}

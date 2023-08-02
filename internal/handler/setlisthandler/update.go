package setlisthandler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/util"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type setlistUpdateReq struct {
	Name           string                 `json:"name" binding:"required"`
	CreatorID      int64                  `json:"creator_id" binding:"required"`
	Deadline       time.Time              `json:"deadline" binding:"required"`
	CreatedEntries []setlistRoleCreateReq `json:"created_entries" binding:"required,dive"`
	UpdatedEntries []setlistRoleUpdateReq `json:"updated_entries" binding:"required,dive"`
	DeletedEntries []int64                `json:"deleted_entries" binding:"required"`
}

type setlistRoleUpdateReq struct {
	ID          int64    `json:"id" binding:"required"`
	SongID      int64    `json:"song_id" binding:"required"`
	Transpose   int16    `json:"transpose"`
	Notes       string   `json:"notes"`
	Arrangement []string `json:"arrangement"`
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
	if err := util.BindModel(ctx, &slReq); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

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
		Deadline:  slReq.Deadline.Local().Truncate((time.Minute)),
	}

	context := ctx.Request.Context()
	updatedSetlist, err := slh.sls.Update(context, setlist, user)

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	createdEntries := make([]domain.SetlistEntry, len(slReq.CreatedEntries))
	updatedEntries := make([]domain.SetlistEntry, len(slReq.UpdatedEntries))

	for idx, entry := range slReq.CreatedEntries {
		jsonArray, _ := json.Marshal(entry.Arrangement)

		createdEntries[idx] = domain.SetlistEntry{
			SongID:      entry.SongID,
			SetlistID:   updatedSetlist.ID,
			Transpose:   entry.Transpose,
			Notes:       entry.Notes,
			Arrangement: datatypes.JSON(jsonArray),
		}
	}

	for idx, entry := range slReq.UpdatedEntries {
		jsonArray, _ := json.Marshal(entry.Arrangement)

		updatedEntries[idx] = domain.SetlistEntry{
			ID:          entry.ID,
			SongID:      entry.SongID,
			SetlistID:   updatedSetlist.ID,
			Transpose:   entry.Transpose,
			Notes:       entry.Notes,
			Arrangement: datatypes.JSON(jsonArray),
		}
	}

	if err := slh.sles.StoreBatch(context, &createdEntries, user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	if err := slh.sles.UpdateBatch(context, &updatedEntries, user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	if err := slh.sles.RemoveBatch(context, setlist, slReq.DeletedEntries, user); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})

		return
	}

	entries := append(createdEntries, updatedEntries...)

	response := setlistResponse{
		updatedSetlist,
		&entries,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"setlist": response,
	})
}

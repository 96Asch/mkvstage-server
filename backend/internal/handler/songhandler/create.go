package songhandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type songCreateReq struct {
	BundleID   int64  `json:"bundle_id" binding:"required"`
	Title      string `json:"title" binding:"required,lte=255"`
	Subtitle   string `json:"subtitle" binding:"lte=255"`
	Key        string `json:"key" binding:"required"`
	Bpm        uint   `json:"bpm" binding:"required"`
	ChordSheet string `json:"chord_sheet" binding:"required"`
}

func (sh songHandler) Create(ctx *gin.Context) {
	val, exists := ctx.Get("user")
	if !exists {
		newErr := domain.NewInternalErr()
		ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})

		return
	}

	var sReq songCreateReq
	if err := ctx.BindJSON(&sReq); err != nil {
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

	context := ctx.Request.Context()
	song := &domain.Song{
		BundleID:   sReq.BundleID,
		CreatorID:  user.ID,
		Title:      sReq.Title,
		Subtitle:   sReq.Subtitle,
		Key:        sReq.Key,
		Bpm:        sReq.Bpm,
		ChordSheet: datatypes.JSON([]byte(sReq.ChordSheet)),
	}

	err := sh.ss.Store(context, song, user)
	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"song": song})
}

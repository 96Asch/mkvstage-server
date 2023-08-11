package setlistrolehandler

import (
	"net/http"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/util"
	"github.com/gin-gonic/gin"
)

type setlistRolePair struct {
	SetlistID  int64 `json:"setlist_id" binding:"required"`
	UserRoleID int64 `json:"userrole_id" binding:"required"`
}

type setlistRoleCreateReq struct {
	UserRoles []setlistRolePair `json:"userroles" binding:"required,dive"`
}

func (srh setlistRoleHandler) Create(ctx *gin.Context) {
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

	var slrReq setlistRoleCreateReq
	if err := util.BindModel(ctx, &slrReq); err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	setlistRoles := make([]domain.SetlistRole, len(slrReq.UserRoles))

	for idx, userrole := range slrReq.UserRoles {
		setlistRoles[idx].SetlistID = userrole.SetlistID
		setlistRoles[idx].UserRoleID = userrole.UserRoleID
	}

	context := ctx.Request.Context()
	err := srh.slrs.Store(context, &setlistRoles, user)

	if err != nil {
		ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"userroles": setlistRoles,
	})
}

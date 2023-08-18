package bundlehandler

import (
	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type bundleHandler struct {
	bs domain.BundleService
}

func Initialize(rg *gin.RouterGroup, bs domain.BundleService, mwh domain.MiddlewareHandler) {
	bundlehandler := &bundleHandler{bs: bs}

	bundle := rg.Group("bundles")
	bundle.GET(":id", bundlehandler.GetByID)
	bundle.GET("", bundlehandler.GetAll)
	bundle.POST("create", mwh.AuthenticateUser(), bundlehandler.Create)
	bundle.DELETE(":id/delete", mwh.AuthenticateUser(), bundlehandler.Delete)
	bundle.PUT(":id/update", mwh.AuthenticateUser(), bundlehandler.UpdateByID)
}

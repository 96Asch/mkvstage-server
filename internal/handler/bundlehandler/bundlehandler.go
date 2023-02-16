package bundlehandler

import (
	"log"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type bundleHandler struct {
	bs domain.BundleService
}

func Initialize(rg *gin.RouterGroup, bs domain.BundleService, mwh domain.MiddlewareHandler) {
	log.Println("Setting up bundle handlers")
	bh := &bundleHandler{
		bs: bs,
	}

	bundle := rg.Group("bundles")
	bundle.GET(":id", bh.Get)
	bundle.GET("", bh.GetAll)
	bundle.POST("create", mwh.AuthenticateUser(), bh.Create)
	bundle.DELETE(":id/delete", mwh.AuthenticateUser(), bh.Delete)
	bundle.PUT(":id/update", mwh.AuthenticateUser(), bh.UpdateByID)

}

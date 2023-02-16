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
	bundle.POST("create", mwh.AuthenticateUser(), bh.Create)

}

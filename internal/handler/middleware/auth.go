package middleware

import (
	"strings"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

type tokenHeader struct {
	Access string `header:"Authorization"`
}

type GinMiddlewareHandler struct {
	TS domain.TokenService
}

func NewGinMiddlewareHandler(ts domain.TokenService) *GinMiddlewareHandler {
	return &GinMiddlewareHandler{
		TS: ts,
	}
}

func (gmh GinMiddlewareHandler) AuthenticateUser() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		header := tokenHeader{}

		if err := ctx.BindHeader(&header); err != nil {
			newErr := domain.NewBadRequestErr(err.Error())
			ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})
			ctx.Abort()
			return
		}

		split := strings.Split(header.Access, " ")

		if len(split) != 2 {
			newErr := domain.NewBadRequestErr("incorrect number of token arguments")
			ctx.JSON(domain.Status(newErr), gin.H{"error": newErr})
			ctx.Abort()
			return
		}

		context := ctx.Request.Context()
		user, err := gmh.TS.ExtractUser(context, split[1])
		if err != nil {
			ctx.JSON(domain.Status(err), gin.H{"error": err})
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}

}

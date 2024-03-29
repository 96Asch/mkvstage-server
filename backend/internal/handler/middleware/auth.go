package middleware

import (
	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type tokenHeader struct {
	Access string `header:"Authorization"`
}

func (gmh ginMiddlewareHandler) AuthenticateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := tokenHeader{}

		if err := ctx.BindHeader(&header); err != nil {
			newErr := domain.NewBadRequestErr(err.Error())
			ctx.JSON(domain.Status(newErr), gin.H{"error": newErr.Error()})
			ctx.Abort()

			return
		}

		context := ctx.Request.Context()

		email, err := gmh.TS.ExtractEmail(context, header.Access)
		if err != nil {
			ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})
			ctx.Abort()

			return
		}

		user, err := gmh.US.FetchByEmail(context, email)
		if err != nil {
			ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})
			ctx.Abort()

			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}

func (gmh ginMiddlewareHandler) JWTExtractEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := tokenHeader{}

		if err := ctx.BindHeader(&header); err != nil {
			newErr := domain.NewBadRequestErr(err.Error())
			ctx.JSON(domain.Status(newErr), gin.H{"error": newErr.Error()})
			ctx.Abort()

			return
		}

		context := ctx.Request.Context()

		email, err := gmh.TS.ExtractEmail(context, header.Access)
		if err != nil {
			ctx.JSON(domain.Status(err), gin.H{"error": err.Error()})
			ctx.Abort()

			return
		}

		ctx.Set("email", email)
		ctx.Next()
	}
}

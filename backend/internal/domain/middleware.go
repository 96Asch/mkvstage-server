package domain

import "github.com/gin-gonic/gin"

type MiddlewareHandler interface {
	AuthenticateUser() gin.HandlerFunc
}

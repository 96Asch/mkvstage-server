package mocks

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockMiddlewareHandler struct {
	mock.Mock
}

func (m MockMiddlewareHandler) AuthenticateUser() gin.HandlerFunc {
	ret := m.Called()

	var r0 gin.HandlerFunc
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(gin.HandlerFunc)
	}

	return r0
}

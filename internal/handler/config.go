package handler

import (
	"github.com/gin-gonic/gin"
)

type Config struct {
	Router *gin.Engine
}

func (cfg *Config) New() *Config {
	return &Config{
		Router: gin.Default(),
	}
}

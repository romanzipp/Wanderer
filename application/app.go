package application

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Router *gin.Engine
	Env    Env
}

type Env struct {
	Password        string
	SessionLifetime int
}

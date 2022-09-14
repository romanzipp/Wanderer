package web

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func CreateServerController(c *gin.Context, db *gorm.DB) {
	c.HTML(http.StatusOK, "servers-create", gin.H{
		"title": "Create server",
	})
}

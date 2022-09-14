package web

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func WebIndexController(c *gin.Context, db *gorm.DB) {
	c.HTML(http.StatusOK, "index", gin.H{
		"title": "Main website",
	})
}

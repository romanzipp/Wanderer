package web

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func ShowAuthFormController(c *gin.Context, db *gorm.DB) {
	c.HTML(http.StatusOK, "auth", gin.H{
		"title": "Authenticate",
	})
}

func SubmitAuthController(c *gin.Context, db *gorm.DB) {
	c.SetCookie("token", "ok", 3000, "/", "", false, true)
	c.Redirect(302, "/")
}

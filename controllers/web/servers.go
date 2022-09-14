package web

import (
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/models"
	"gorm.io/gorm"
	"net/http"
)

func ListServersController(c *gin.Context, db *gorm.DB) {
	db.AutoMigrate(&models.Server{})

	var servers []models.Server
	db.Find(&servers)

	c.HTML(http.StatusOK, "servers", gin.H{
		"title":   "Server",
		"servers": servers,
	})
}
func ShowCreateServerController(c *gin.Context, db *gorm.DB) {
	c.HTML(http.StatusOK, "servers-create", gin.H{
		"title": "Create server",
	})
}

func CreateServerController(c *gin.Context, db *gorm.DB) {
	db.Create(&models.Server{
		Name:                 c.PostForm("name"),
		Address:              c.PostForm("address"),
		Port:                 c.PostForm("port"),
		CfAccessClientId:     c.PostForm("cf_access_client_id"),
		CfAccessClientSecret: c.PostForm("cf_access_client_secret"),
	})

	c.Redirect(302, "/servers")
}

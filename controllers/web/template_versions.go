package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/models"
	"gorm.io/gorm"
	"strconv"
)

func CreateVersionController(c *gin.Context, db *gorm.DB, templateID string) {
	templateIDConv, _ := strconv.ParseInt(templateID, 10, 64)

	db.Create(&models.TemplateVersion{
		Selector:    c.PostForm("selector"),
		LastVersion: c.PostForm("version"),
		TemplateID:  int(templateIDConv),
	})

	c.Redirect(302, fmt.Sprintf("/templates/%s?success=Version+created", templateID))
}

func DeleteVersionController(c *gin.Context, db *gorm.DB, versionID string) {
	versionIDConv, _ := strconv.ParseInt(versionID, 10, 64)
	version := &models.TemplateVersion{}
	version.ID = uint(versionIDConv)

	db.Delete(version)

	c.Redirect(302, fmt.Sprintf("/templates?success=Version+deleted"))
}

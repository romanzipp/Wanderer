package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/application"
	"github.com/romanzipp/wanderer/models"
	"strconv"
)

func CreateVersionController(c *gin.Context, app *application.App, templateID string) {
	templateIDConv, _ := strconv.ParseInt(templateID, 10, 64)

	app.DB.Create(&models.TemplateVersion{
		Selector:    c.PostForm("selector"),
		LastVersion: c.PostForm("version"),
		TemplateID:  int(templateIDConv),
	})

	c.Redirect(302, fmt.Sprintf("/templates/%s?success=Version+created", templateID))
}

func DeleteVersionController(c *gin.Context, app *application.App, versionID string) {
	versionIDConv, _ := strconv.ParseInt(versionID, 10, 64)
	version := &models.TemplateVersion{}
	version.ID = uint(versionIDConv)

	app.DB.Delete(version)

	c.Redirect(302, fmt.Sprintf("/templates?success=Version+deleted"))
}

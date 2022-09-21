package api

import (
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/application"
)

type IndexResponse struct {
	Status string
}

func IndexController(c *gin.Context, app *application.App) {
	c.JSON(200, IndexResponse{"Ok"})
}

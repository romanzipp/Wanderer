package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/application"
	"github.com/romanzipp/wanderer/controllers/api"
	"github.com/romanzipp/wanderer/models"
	"gorm.io/gorm"
	"time"
)

type JsonError struct {
	Message string `json:"message"`
}

func InitApiRoutes(app *application.App) {
	authed := app.Router.Group("/api")
	authed.Use(ApiTokenAuth(app.DB))

	authed.GET("/", func(c *gin.Context) {
		api.IndexController(c, app)
	})

	authed.POST("/deploy", func(c *gin.Context) {
		api.DeployController(c, app)
	})
}

func ApiTokenAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")

		if len(tokenStr) == 0 {
			c.AbortWithStatusJSON(401, JsonError{"Missing token"})
			return
		}

		var token models.Token
		db.Where("token = ?", tokenStr).Find(&token)

		if token.ID == 0 {
			c.AbortWithStatusJSON(401, JsonError{fmt.Sprintf("Invalid token: %s", tokenStr)})
			return
		}

		token.LastUsedAt = time.Now()
		db.Save(token)
		c.Next()
	}
}

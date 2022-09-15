package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/controllers/api"
	"github.com/romanzipp/wanderer/models"
	"gorm.io/gorm"
)

type JsonError struct {
	Message string
}

func InitApiRoutes(router *gin.Engine, db *gorm.DB) {
	authed := router.Group("/api")
	authed.Use(ApiTokenAuth(db))

	authed.GET("/", func(c *gin.Context) {
		api.IndexController(c, db)
	})
}

func ApiTokenAuth(db *gorm.DB) gin.HandlerFunc {
	abort := func(c *gin.Context, code int, message string) {
		c.JSON(code, JsonError{message})
	}

	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")

		if len(tokenStr) == 0 {
			abort(c, 401, "Missing token")
			return
		}

		var token models.Token
		db.Where("token = ?", tokenStr).Find(&token)

		if token.ID == 0 {
			abort(c, 401, "Invalid token")
			return
		}
	}
}

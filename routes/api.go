package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/romanzipp/wanderer/controllers/api"
	"github.com/romanzipp/wanderer/models"
	"gorm.io/gorm"
	"time"
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

	authed.POST("/deploy", func(c *gin.Context) {
		api.DeployController(c, db)
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
			c.AbortWithStatusJSON(401, JsonError{"Invalid token"})
			return
		}

		token.LastUsedAt = time.Now()
		db.Save(token)
		c.Next()
	}
}

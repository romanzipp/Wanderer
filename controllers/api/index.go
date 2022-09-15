package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IndexResponse struct {
	Status string
}

func IndexController(c *gin.Context, db *gorm.DB) {
	c.JSON(200, IndexResponse{"Ok"})
}

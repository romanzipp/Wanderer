package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/romanzipp/wanderer/models"
	"gorm.io/gorm"
	"net/http"
)

func ApiController(c *gin.Context, db *gorm.DB) {
	var tokens []models.Token
	db.Find(&tokens)

	token := c.Query("token")
	c.HTML(http.StatusOK, "api", gin.H{
		"title":        "API",
		"nav":          "api",
		"createdToken": token,
		"tokens":       tokens,
	})
}

func IssueApiTokenController(c *gin.Context, db *gorm.DB) {
	token := uuid.NewString()
	db.Create(&models.Token{
		Name:  c.PostForm("name"),
		Token: token,
	})

	c.Redirect(302, fmt.Sprintf("/tokens?token=%s", token))
}

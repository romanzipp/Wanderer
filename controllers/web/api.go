package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/romanzipp/wanderer/application"
	"github.com/romanzipp/wanderer/models"
	"net/http"
)

func ApiController(c *gin.Context, app *application.App) {
	var tokens []models.Token
	app.DB.Find(&tokens)

	token := c.Query("token")
	c.HTML(http.StatusOK, "api", gin.H{
		"title":        "API",
		"nav":          "api",
		"createdToken": token,
		"tokens":       tokens,
	})
}

func IssueApiTokenController(c *gin.Context, app *application.App) {
	token := uuid.NewString()
	app.DB.Create(&models.Token{
		Name:  c.PostForm("name"),
		Token: token,
	})

	c.Redirect(302, fmt.Sprintf("/tokens?token=%s&success=Token+issued", token))
}

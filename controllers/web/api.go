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
		"success":      c.Query("success"),
		"currentUrl":   c.Request.URL.Path,
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

func DeleteApiTokenController(c *gin.Context, app *application.App, tokenID string) {
	var token models.Token
	app.DB.Where("id = ?", tokenID).First(&token)

	if token.ID == 0 {
		c.Redirect(302, "/tokens?error=Token+not+found")
		return
	}

	app.DB.Delete(&token)

	c.Redirect(302, "/tokens?success=Token+deleted")
}

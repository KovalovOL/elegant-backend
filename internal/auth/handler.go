package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(s *AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	state, url := h.service.StartGoogleLogin()
	c.SetCookie("oauth_state", state, 7*60, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	storedState, err := c.Cookie("oauth_state")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing state"})
		return
	}

	if c.Query("state") != storedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	code := c.Query("code")
	_, token, err := h.service.HandleGoogleCallback(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 24*60*60, "/", "", false, true)

	c.Redirect(http.StatusSeeOther, os.Getenv("FRONTEND_URL")+"/dashboard")
}

func (h *AuthHandler) Me(c *gin.Context) {
	id := c.GetString("user_id")
	name := c.GetString("user_name")
	email := c.GetString("user_email")
	c.JSON(http.StatusOK, gin.H{
		"id":    id,
		"name":  name,
		"email": email,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {

	c.SetCookie("token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out and user deleted"})
}

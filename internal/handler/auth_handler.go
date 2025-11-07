package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"app/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
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
	user, token, err := h.service.HandleGoogleCallback(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 24*60*60, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"name": user.Name})
}

// 3️⃣ /me
func (h *AuthHandler) Me(c *gin.Context) {
	name := c.GetString("user_name")
	c.JSON(http.StatusOK, gin.H{"name": name})
}

package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GoogleLogin(c *gin.Context) {
	state, url := h.service.StartGoogleLogin()
	c.SetCookie("oauth_state", state, 7*60, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) GoogleCallback(c *gin.Context) {
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

	frontendURL := "http://localhost:5173" //os.Getenv("FRONTEND_URL")
	c.Redirect(http.StatusSeeOther, frontendURL+"/dashboard")
}

func (h *Handler) Me(c *gin.Context) {
	id := c.GetString("user_id")
	name := c.GetString("user_name")
	email := c.GetString("user_email")
	c.JSON(http.StatusOK, gin.H{
		"id":    id,
		"name":  name,
		"email": email,
	})
}

func (h *Handler) Logout(c *gin.Context) {

	c.SetCookie("token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out and user deleted"})
}

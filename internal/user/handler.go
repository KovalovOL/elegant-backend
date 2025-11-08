package user

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) DeleteCurrentUser(c *gin.Context) {
	ctx := c.Request.Context()
	err := h.service.DeleteCurrentUser(ctx, c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "user deleted"})
}

func (h *Handler) UpdateCurrentUser(c *gin.Context) {
	var user CreateUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	ctx := c.Request.Context()
	err := h.service.UpdateCurrentUser(ctx, c, user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "user updated"})
}

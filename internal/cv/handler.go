package cv

import (
	"github.com/gin-gonic/gin"
	"strconv"
)


type CVHandler struct {
	service *CVService
}

func NewCVHandler(s *CVService) *CVHandler {
	return &CVHandler{service: s}
}


func (h *CVHandler) GetAllCVsByCurrentUser(c *gin.Context) {
	ctx := c.Request.Context()
	cvs, err := h.service.GetAllCVsByCurrentUser(ctx, c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, cvs)
}

func (h *CVHandler) GetCVByID(c *gin.Context) {
	cvID, err := strconv.Atoi(c.Param("cv_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid cv_id"})
		return
	}

	ctx := c.Request.Context()
	cv, err := h.service.GetCVByID(ctx, c, cvID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, cv)
}

func (h *CVHandler) CreateCV(c *gin.Context) {
	var input CreateCVInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	ctx := c.Request.Context()
	id, err := h.service.CreateCV(ctx, c, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"cv_id": id})
}

func (h *CVHandler) DeleteCV(c *gin.Context) {
	cvID, err := strconv.Atoi(c.Param("cv_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid cv_id"})
		return
	}

	ctx := c.Request.Context()
	err = h.service.DeleteCV(ctx, c, cvID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "cv deleted"})
}

func (h *CVHandler) UpdateCV(c *gin.Context) {
	cvID, err := strconv.Atoi(c.Param("cv_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid cv_id"})
		return
	}

	var input CreateCVInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	ctx := c.Request.Context()
	err = h.service.UpdateCV(ctx, c, cvID, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "cv updated"})
}
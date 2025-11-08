package tag

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type TagHandler struct {
	service *TagService
}

func NewTagHandler(service *TagService) *TagHandler {
	return &TagHandler{service}
}

func (h *TagHandler) GetAllTags(c *gin.Context){
	ctx := c.Request.Context()
	tags, err := h.service.GetAllTags(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, tags)
}

func (h *TagHandler) GetTagById(c *gin.Context) {
	tagID, err := strconv.Atoi(c.Param("tag_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid cv_id"})
		return
	}

	ctx := c.Request.Context()
	tag, err := h.service.GetTagById(ctx, tagID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, tag)
}




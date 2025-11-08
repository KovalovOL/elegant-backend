package test

import (
	"github.com/gin-gonic/gin"
	"strings"
	"strconv"
)

type TestHandler struct {
	service *TestService
}

func NewTestHandler(service *TestService) *TestHandler {
	return &TestHandler{service}
}

func (h *TestHandler) CreateTest(c *gin.Context) {
	ctx := c.Request.Context()
	var req CreateTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if len(req.TagIDs) == 0 {
		c.JSON(400, gin.H{"error": "tag_ids cannot be empty"})
		return
	}
	newTest := CreateTest{
		Title: req.Title,
		TimeLimit: req.TimeLimit,
		Type: req.Type,
		Tasks: req.Tasks,
	}
	id, err := h.service.CreateTest(ctx, newTest, req.TagIDs)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"test_id": id})
}

func (h *TestHandler) GetTestsByTags(c *gin.Context) {
    tagIDs := parseTagIDs(c.Query("tags"))
    
    var tests []Test
    var err error
    
    if len(tagIDs) > 0 {
        tests, err = h.service.GetTestsByTags(c.Request.Context(), tagIDs)
    } else {
        tests, err = h.service.GetAllTests(c.Request.Context())
    }
    
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, tests)
}

func parseTagIDs(tagStr string) []int {
    if tagStr == "" {
        return nil
    }
    
    tagIDs := []int{}
    tagStrs := strings.Split(tagStr, ",")
    
    for _, ts := range tagStrs {
        if id, err := strconv.Atoi(strings.TrimSpace(ts)); err == nil {
            tagIDs = append(tagIDs, id)
        }
    }
    
    return tagIDs
}
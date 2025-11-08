package user

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) DeleteCurrentUser(ctx context.Context, c *gin.Context) error {
	id, _ := strconv.Atoi(c.GetString("user_id"))
	err := s.repo.DeleteUserById(c, id)
	if err != nil {
		return err
	}
	c.Set("user_id", nil)
	c.Set("user_name", nil)
	c.Set("user_email", nil)
	return nil
}

func (s *Service) UpdateCurrentUser(ctx context.Context, c *gin.Context, newUser CreateUser) error {
	id, _ := strconv.Atoi(c.GetString("user_id"))
	err := s.repo.UpdateUserById(ctx, id, newUser)
	if err != nil {
		return err
	}
	return nil
}

package cv

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)


type CVService struct {
	repo *CVRepository
}

func NewCVService(r *CVRepository) *CVService {
	return &CVService{repo: r}
}


func (s *CVService) GetAllCVsByCurrentUser(ctx context.Context, c *gin.Context) ([]CV, error) {
	userID, _ := strconv.Atoi(c.GetString("user_id"))
	return s.repo.GetByUserID(ctx, userID)
}

func (s *CVService) GetCVByID(ctx context.Context, c *gin.Context, cvID int) (*CV, error) {
	current_user_id, _ := strconv.Atoi(c.GetString("user_id"))
	cv, err := s.repo.GetByID(ctx, cvID)
	if err != nil {
		return nil, err
	}
	if cv.UserID != current_user_id {
		return nil, gin.Error{Err: err, Type: gin.ErrorTypePublic, Meta: "unauthorized access to CV"}
	}

	return cv, nil
}

func (s *CVService) CreateCV(ctx context.Context, c *gin.Context, input *CreateCVInput) (int, error) {
	current_user_id, _ := strconv.Atoi(c.GetString("user_id"))
	cv := &CreateCV{
		CreateCVInput: *input,
		UserID:        current_user_id,
	}
	id, err := s.repo.Create(ctx, cv)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *CVService) DeleteCV(ctx context.Context, c *gin.Context, cvID int) error {
	current_user_id, _ := strconv.Atoi(c.GetString("user_id"))
	cv, err := s.repo.GetByID(ctx, cvID)
	if err != nil {
		return err
	}
	if cv.UserID != current_user_id {
		return gin.Error{Err: err, Type: gin.ErrorTypePublic, Meta: "unauthorized access to CV"}
	}

	return s.repo.Delete(ctx, cvID)
}

func (s *CVService) UpdateCV(ctx context.Context, c *gin.Context, cvID int, input *CreateCVInput) error {
	current_user_id, _ := strconv.Atoi(c.GetString("user_id"))
	cv, err := s.repo.GetByID(ctx, cvID)
	if err != nil {
		return err
	}
	if cv.UserID != current_user_id {
		return gin.Error{Err: err, Type: gin.ErrorTypePublic, Meta: "unauthorized access to CV"}
	}

	return s.repo.Update(ctx , &CV{
		CreateCV: CreateCV{
			CreateCVInput: *input,
			UserID: current_user_id,
		},
		CVID: cvID,
	})
}
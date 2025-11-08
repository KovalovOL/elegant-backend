package test

import "context"

type TestService struct {
	repo *TestRepository
}

func NewTestService(repo *TestRepository) *TestService {
	return &TestService{repo}
}

func (s *TestService) CreateTest(ctx context.Context, newTest CreateTest, tagIDs []int) (int, error){
	return s.repo.CreateTestFull(ctx, newTest, tagIDs)
}

func (s *TestService) GetAllTests(ctx context.Context) ([]Test, error) {
	return s.repo.GetAllTests(ctx)
}

func (s *TestService) GetTestById(ctx context.Context, testID int) (*Test, error) {
	return s.repo.GetTestByID(ctx, testID)
}

func (s *TestService) GetTestsByTags(ctx context.Context, tagIDs []int) ([]Test, error) {
	return s.repo.GetTestsByTags(ctx, tagIDs)
}

func (s *TestService) DeleteTestById(ctx context.Context, testID int) error {
	return s.repo.DeleteTest(ctx, testID)
}

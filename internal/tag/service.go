package tag

import (
	"context"
)


type TagService struct {
	repo *TagRepository
}

func NewTagService(repo *TagRepository) *TagService {
	return &TagService{repo}
}


func (s *TagService) GetAllTags(ctx context.Context) ([]Tag, error) {
	tagPtrs, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	tags := make([]Tag, len(tagPtrs))
	for i, t := range tagPtrs {
		tags[i] = *t
	}

	return tags, nil
}

func (s *TagService) GetTagById(ctx context.Context, tag_id int) (*Tag, error) {
	tag, err := s.repo.GetByID(ctx, tag_id)
	if err != nil {
		return nil, err
	}
	return tag, nil
}
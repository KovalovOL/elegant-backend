package tag

import (
	"context"
	"database/sql"
	"fmt"
)

type TagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(ctx context.Context, input *CreateTag) (int, error) {
	query := `
		INSERT INTO tags (name)
		VALUES ($1)
		RETURNING tag_id
	`
	var id int
	err := r.db.QueryRowContext(ctx, query, input.Name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create tag: %w", err)
	}

	return id, nil
}

func (r *TagRepository) GetByID(ctx context.Context, id int) (*Tag, error) {
	query := `
		SELECT tag_id, name
		FROM tags
		WHERE tag_id = $1
	`

	var t Tag
	err := r.db.QueryRowContext(ctx, query, id).Scan(&t.TagID, &t.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	return &t, nil
}

func (r *TagRepository) GetAll(ctx context.Context) ([]*Tag, error) {
	query := `
		SELECT tag_id, name
		FROM tags
		ORDER BY tag_id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	defer rows.Close()

	var tags []*Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.TagID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, &t)
	}

	return tags, nil
}

func (r *TagRepository) Update(ctx context.Context, id int, input *CreateTag) error {
	query := `
		UPDATE tags
		SET name = $1
		WHERE tag_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, input.Name, id)
	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}

	return nil
}

func (r *TagRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM tags
		WHERE tag_id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

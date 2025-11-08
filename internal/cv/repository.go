package cv

import (
	"context"
	"database/sql"
)

type CVRepository struct {
	db *sql.DB
}

func NewCVRepository(db *sql.DB) *CVRepository {
	return &CVRepository{db: db}
}

func (r *CVRepository) Create(ctx context.Context, input *CreateCV) (int, error) {
	query := `
		INSERT INTO cvs (user_id, title, position, summary, skills, experience, education)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING cv_id
	`
	var id int
	err := r.db.QueryRowContext(ctx, query,
		input.UserID,
		input.Title,
		input.Position,
		input.Summary,
		input.Skills,
		input.Experience,
		input.Education,
	).Scan(&id)
	return id, err
}

func (r *CVRepository) GetByID(ctx context.Context, id int) (*CV, error) {
	query := `
		SELECT cv_id, user_id, title, position, summary, skills, experience, education
		FROM cvs
		WHERE cv_id = $1
	`
	var cv CV
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cv.CVID,
		&cv.UserID,
		&cv.Title,
		&cv.Position,
		&cv.Summary,
		&cv.Skills,
		&cv.Experience,
		&cv.Education,
	)
	if err != nil {
		return nil, err
	}
	return &cv, nil
}

func (r *CVRepository) GetByUserID(ctx context.Context, userID int) ([]CV, error) {
	query := `
		SELECT cv_id, user_id, title, position, summary, skills, experience, education
		FROM cvs
		WHERE user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cvs []CV
	for rows.Next() {
		var cv CV
		if err := rows.Scan(
			&cv.CVID,
			&cv.UserID,
			&cv.Title,
			&cv.Position,
			&cv.Summary,
			&cv.Skills,
			&cv.Experience,
			&cv.Education,
		); err != nil {
			return nil, err
		}
		cvs = append(cvs, cv)
	}
	return cvs, nil
}

func (r *CVRepository) Update(ctx context.Context, cv *CV) error {
	query := `
		UPDATE cvs
		SET title = $1,
		    position = $2,
		    summary = $3,
		    skills = $4,
		    experience = $5,
		    education = $6
		WHERE cv_id = $7 AND user_id = $8
	`
	_, err := r.db.ExecContext(ctx, query,
		cv.Title,
		cv.Position,
		cv.Summary,
		cv.Skills,
		cv.Experience,
		cv.Education,
		cv.CVID,
		cv.UserID,
	)
	return err
}

func (r *CVRepository) Delete(ctx context.Context, cvID int) error {
	query := `DELETE FROM cvs WHERE cv_id = $1`
	_, err := r.db.ExecContext(ctx, query, cvID)
	return err
}

package user

import (
	"context"
	"database/sql"
)

// type repository interface{}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user CreateUser) (int, error) {
	var userID int
	err := r.db.QueryRow(
		`INSERT INTO users (email, name, github_url, linkedin_url, bio) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, user.Name, user.GitHubUrl, user.LikedinUrl, user.Bio,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.QueryRow(
		`SELECT id, email, name, github_url, linkedin_url, bio 
		 FROM users WHERE email = $1`, email,
	).Scan(&user.ID, &user.Email, &user.Name, &user.GitHubUrl, &user.LikedinUrl, &user.Bio)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserById(ctx context.Context, id int) (*User, error) {
	var user User
	err := r.db.QueryRow(
		`SELECT id, email, name, github_url, linkedin_url, bio 
		 FROM users WHERE id = $1`, id,
	).Scan(&user.ID, &user.Email, &user.Name, &user.GitHubUrl, &user.LikedinUrl, &user.Bio)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) DeleteUserById(ctx context.Context, id int) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *Repository) UpdateUserById(ctx context.Context, id int, newUser CreateUser) error {
	_, err := r.db.Exec(
		`UPDATE users SET name = $1, github_url = $2, linkedin_url = $3, bio = $4 
		 WHERE id = $5`,
		newUser.Name, newUser.GitHubUrl, newUser.LikedinUrl, newUser.Bio, id,
	)
	return err
}

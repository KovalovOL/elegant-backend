package auth

import (
	"app/internal/user"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
)

type AuthService struct {
	oauth     *GoogleOAuth
	jwt       *JWTManager
	user_repo *user.UserRepository
}

func NewAuthService(oauth *GoogleOAuth, jwt *JWTManager, user_repo *user.UserRepository) *AuthService {
	return &AuthService{oauth: oauth, jwt: jwt, user_repo: user_repo}
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *AuthService) StartGoogleLogin() (state, url string) {
	state = generateState()
	url = s.oauth.GetAuthURL(state)
	return state, url
}

func (s *AuthService) HandleGoogleCallback(ctx context.Context, code string) (*user.User, string, error) {
	token, err := s.oauth.ExchangeCode(ctx, code)
	if err != nil {
		return nil, "", err
	}

	userByJWT, err := s.oauth.GetUserInfo(ctx, token)
	if err != nil {
		return nil, "", err
	}

	// Check if user exists in DB
	existingUser, err := s.user_repo.GetUserByEmail(ctx, userByJWT.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, "", err
	}

	// If not exists — create
	if existingUser == nil {
		newUser := user.CreateUser{
			Email: userByJWT.Email,
			Name:  userByJWT.Name,
		}

		id, err := s.user_repo.CreateUser(ctx, newUser)
		if err != nil {
			return nil, "", err
		}

		existingUser, err = s.user_repo.GetUserById(ctx, id)
		if err != nil {
			return nil, "", err
		}
	}

	if existingUser == nil {
		return nil, "", fmt.Errorf("failed to find or create user")
	}

	jwt, err := s.jwt.Generate(existingUser)
	if err != nil {
		return nil, "", err
	}

	return existingUser, jwt, nil
}

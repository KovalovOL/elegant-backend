package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"app/internal/auth"
	"app/internal/model"
)

type AuthService struct {
	oauth *auth.GoogleOAuth
	jwt   *auth.JWTManager
}

func NewAuthService(oauth *auth.GoogleOAuth, jwt *auth.JWTManager) *AuthService {
	return &AuthService{oauth: oauth, jwt: jwt}
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

func (s *AuthService) HandleGoogleCallback(ctx context.Context, code string) (*model.UserGoogleResp, string, error) {
	token, err := s.oauth.ExchangeCode(ctx, code)
	if err != nil {
		return nil, "", err
	}

	user, err := s.oauth.GetUserInfo(ctx, token)
	if err != nil {
		return nil, "", err
	}

	jwt, err := s.jwt.Generate(user)
	if err != nil {
		return nil, "", err
	}

	return user, jwt, nil
}

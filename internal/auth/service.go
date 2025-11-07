package auth

import (
	"app/internal/user"
	"context"
	"crypto/rand"
	"encoding/hex"
)

type AuthService struct {
	oauth *GoogleOAuth
	jwt   *JWTManager
}

func NewAuthService(oauth *GoogleOAuth, jwt *JWTManager) *AuthService {
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

func (s *AuthService) HandleGoogleCallback(ctx context.Context, code string) (*user.UserGoogleResp, string, error) {
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

package auth_manager

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"seat-reservation/jwt"
	"seat-reservation/pkg/manager/auth/repository"
	"seat-reservation/pkg/manager/port"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	AuthToken    string
	RefreshToken string
}

type authManager struct {
	authRepository repository.Repository
	userPort       port.UserPort
	jwtConfig      jwt.Config
}

func NewAuthManager(authRepository repository.Repository, userPort port.UserPort, jwtConfig jwt.Config) AuthManager {
	return &authManager{
		authRepository,
		userPort,
		jwtConfig,
	}
}

func (m *authManager) SignupOrLogin(ctx context.Context, req SignupOrLoginRequest) (*SignupOrLoginResponse, error) {
	// create user
	createUserReq := port.CreateOrGetUserRequest{
		Username: req.UserName,
		Password: req.Password,
	}

	user, err := m.userPort.CreateOrGetUser(ctx, createUserReq)
	if err != nil {
		return nil, err
	}

	if user.Password != req.Password {
		return nil, fmt.Errorf("invalid password")
	}

	// create token
	token, err := m.issueTokens(ctx, user.ID, nil)
	if err != nil {
		return nil, err
	}

	// return
	return &SignupOrLoginResponse{
		Id:           user.ID,
		UserName:     user.Username,
		AuthToken:    token.AuthToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (h *authManager) issueTokens(ctx context.Context, userId uuid.UUID, oldToken *string) (Token, error) {
	tokens, err := h.createTokens(ctx, userId)
	if err != nil {
		return tokens, err
	}

	if oldToken != nil {
		err = h.authRepository.DeleteRefreshToken(ctx, userId)
		if err != nil {
			return tokens, err
		}
	}
	input := repository.RefreshToken{
		UserId:               userId,
		RefreshToken:         tokens.RefreshToken,
		RefreshTokenExpiryAt: time.Now().Add(h.jwtConfig.RefreshTokenExpiryDuration),
	}
	err = h.authRepository.AddRefreshToken(ctx, &input)
	return tokens, err
}

func (h *authManager) createTokens(_ context.Context, userId uuid.UUID) (Token, error) {
	authToken, err := jwt.IssueToken(h.jwtConfig, userId, map[string]interface{}{
		"id": userId,
	})
	if err != nil {
		return Token{}, fmt.Errorf("could not issue JWT for %v, err: %+v", userId, err)
	}

	refreshToken, err := GenerateRandomHex(128)
	if err != nil {
		return Token{}, fmt.Errorf("could not generate refresh token for %v", userId)
	}

	return Token{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}, nil
}

func GenerateRandomHex(length int) (string, error) {
	byteLength := length / 2
	randBytes := make([]byte, byteLength)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randBytes), nil
}

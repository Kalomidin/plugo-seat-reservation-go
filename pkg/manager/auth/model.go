package auth_manager

import (
	"context"

	"github.com/google/uuid"
)

type SignupOrLoginRequest struct {
	UserName string
	Password string
}

type SignupOrLoginResponse struct {
	Id           uuid.UUID
	UserName     string
	AuthToken    string
	RefreshToken string
}

type AuthManager interface {
	SignupOrLogin(ctx context.Context, req SignupOrLoginRequest) (*SignupOrLoginResponse, error)
}

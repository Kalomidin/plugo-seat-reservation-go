package repository

import (
	"context"
	"seat-reservation/common"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	RefreshToken         string
	UserId               uuid.UUID
	RefreshTokenExpiryAt time.Time
	common.CreatedDeleted
}

type Repository interface {
	AddRefreshToken(ctx context.Context, deviceRefreshToken *RefreshToken) error
	DeleteRefreshToken(ctx context.Context, userId uuid.UUID) error
	GetRefreshToken(ctx context.Context, refreshToken string) (*RefreshToken, error)
	Migrate() error
	MigrateDown() error
}

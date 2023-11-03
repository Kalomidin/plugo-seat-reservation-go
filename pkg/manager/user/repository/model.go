package repository

import (
	"context"
	"seat-reservation/common"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	Username string
	Password string
	common.CreatedDeleted
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, userId uuid.UUID) (*User, error)
	GetUserByName(ctx context.Context, username string) (*User, error)
	Migrate() error
	MigrateDown() error
}

package user_manager

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Username string
}

type UserManager interface {
	GetUser(ctx context.Context, userID uuid.UUID) (*User, error)
}

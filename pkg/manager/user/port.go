package user_manager

import (
	"context"
	"seat-reservation/pkg/manager/port"
	"seat-reservation/pkg/manager/user/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userPort struct {
	userRepository repository.Repository
}

func NewUserPort(userRepository repository.Repository) port.UserPort {
	return &userPort{
		userRepository,
	}
}

func (p *userPort) CreateOrGetUser(ctx context.Context, req port.CreateOrGetUserRequest) (*port.User, error) {
	user, err := p.userRepository.GetUserByName(ctx, req.Username)
	if err == nil {
		return &port.User{
			ID:       user.ID,
			Username: user.Username,
			Password: user.Password,
		}, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// else create user
	user = &repository.User{
		Username: req.Username,
		Password: req.Password,
	}
	if err = p.userRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return &port.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}, nil
}

func (p *userPort) GetUser(ctx context.Context, userId uuid.UUID) (*port.User, error) {
	user, err := p.userRepository.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &port.User{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

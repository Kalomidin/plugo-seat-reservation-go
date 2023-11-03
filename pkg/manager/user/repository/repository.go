package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	*gorm.DB
}

func NewRepository(ctx context.Context, db *gorm.DB) Repository {
	return &repository{
		db,
	}
}

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	res := r.Create(user)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return fmt.Errorf("more than one row updated")
	}
	return nil
}

func (r *repository) GetUser(ctx context.Context, userId uuid.UUID) (*User, error) {
	var user User
	resp := r.DB.Where("id = ?", userId).First(&user)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &user, nil
}

func (r *repository) GetUserByName(ctx context.Context, username string) (*User, error) {
	var user User
	resp := r.DB.Where("username = ?", username).First(&user)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &user, nil
}

func (r *repository) Migrate() error {
	return r.AutoMigrate(&User{})
}

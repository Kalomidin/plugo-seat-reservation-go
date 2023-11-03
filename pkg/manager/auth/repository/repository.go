package repository

import (
	"context"
	"fmt"
	"log"
	"time"

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

func (r *repository) AddRefreshToken(ctx context.Context, deviceRefreshToken *RefreshToken) error {
	res := r.Create(deviceRefreshToken)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return fmt.Errorf("could not create refresh token")
	}
	return nil
}

func (r *repository) DeleteRefreshToken(ctx context.Context, userID uuid.UUID) error {
	var d RefreshToken
	resp := r.Model(&d).Where(" user_id is null", userID).Update("deleted_at", time.Now())
	if resp.Error != nil {
		return resp.Error
	}
	if resp.RowsAffected > 1 {
		fmt.Println("more than one refresh token is being deleted")
	} else {
		log.Printf("could not find any active refresh token for given id %s\n", userID)
	}
	return nil
}

func (r *repository) GetRefreshToken(ctx context.Context, refreshToken string) (*RefreshToken, error) {
	var d RefreshToken
	resp := r.Model(&d).Where("refresh_token = ? and deleted_at is null", refreshToken).First(&d)

	if resp.Error != nil {
		return nil, resp.Error
	}
	return &d, nil
}

func (r *repository) Migrate() error {
	return r.AutoMigrate(&RefreshToken{})
}

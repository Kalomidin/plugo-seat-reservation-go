package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type seatRepository struct {
	*gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &seatRepository{
		db,
	}
}

func (repo *seatRepository) CreateSeats(ctx context.Context, seats []*Seat) error {
	res := repo.CreateInBatches(&seats, len(seats))
	return res.Error
}

func (repo *seatRepository) GetSeat(ctx context.Context, id uuid.UUID) (*Seat, error) {
	var seat Seat
	resp := repo.DB.Where("id = ?", id).First(&seat)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &seat, nil
}

func (repo *seatRepository) UpdateSeat(ctx context.Context, seat *Seat) error {
	res := repo.Save(seat)
	return res.Error
}

func (repo *seatRepository) Migrate() error {
	return repo.AutoMigrate(&Seat{})
}

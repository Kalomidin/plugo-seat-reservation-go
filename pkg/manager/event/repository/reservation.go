package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type reservationRepository struct {
	*gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{
		db,
	}
}

func (repo *reservationRepository) CreateReservation(ctx context.Context, reservation *Reservation) error {
	res := repo.Create(reservation)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (repo *reservationRepository) GetReservationsForEvent(ctx context.Context, eventID uuid.UUID) ([]Reservation, error) {
	var reservations []Reservation = []Reservation{}
	res := repo.Find(&reservations, "event_id = ?", eventID)
	if res.Error != nil {
		return nil, res.Error
	}
	return reservations, nil
}

func (repo *reservationRepository) GetUserReservations(ctx context.Context, userID uuid.UUID) ([]Reservation, error) {
	var reservations []Reservation = []Reservation{}
	res := repo.Find(&reservations, "user_id = ?", userID)
	if res.Error != nil {
		return nil, res.Error
	}
	return reservations, nil
}

func (repo *reservationRepository) GetReservationForSeat(ctx context.Context, seatID uuid.UUID) (*Reservation, error) {
	var reservation Reservation
	res := repo.First(&reservation, "seat_id = ?", seatID)
	if res.Error != nil {
		return nil, res.Error
	}
	return &reservation, nil
}

func (repo *reservationRepository) GetReservationForEventAndUser(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (*Reservation, error) {
	var reservation Reservation
	res := repo.First(&reservation, "event_id = ? AND user_id = ?", eventID, userID)
	if res.Error != nil {
		return nil, res.Error
	}
	return &reservation, nil
}

func (repo *reservationRepository) HandleWithTransaction(ctx context.Context, fn Transaction) error {
	tx := repo.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	err := fn(ctx, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (repo *reservationRepository) DeleteReservation(ctx context.Context, id uuid.UUID) error {
	return repo.Model(&Reservation{}).Where("id = ? and deleted_at is null", id).Update("deleted_at", time.Now()).Error
}

func (repo *reservationRepository) Migrate() error {
	return repo.AutoMigrate(&Reservation{})
}

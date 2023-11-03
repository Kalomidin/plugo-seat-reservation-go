package repository

import (
	"context"

	"gorm.io/gorm"
)

type TransactionCallBack func(ctx context.Context,
	eventRepo EventRepository,
	seatRepo SeatRepository,
	reservationRepo ReservationRepository,
) (bool, error)

type transactionDB struct {
	*gorm.DB
}

func NewTransactionDB(db *gorm.DB) TransactionDB {
	return &transactionDB{
		db,
	}
}

func (transactionDB *transactionDB) HandleWithTransaction(ctx context.Context, fn TransactionCallBack) error {
	tx := transactionDB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	eventRepo := NewEventRepository(tx)
	seatRepo := NewSeatRepository(tx)
	reservationRepo := NewReservationRepository(tx)
	commit, err := fn(ctx, eventRepo, seatRepo, reservationRepo)

	if err != nil || !commit {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

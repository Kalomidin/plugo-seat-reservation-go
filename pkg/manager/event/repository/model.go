package repository

import (
	"context"
	"seat-reservation/common"

	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	Name         string
	CreatorID    uuid.UUID
	Seats        []Seat        `gorm:"foreignKey:EventID;references:ID"`
	Reservations []Reservation `gorm:"foreignKey:EventID;references:ID"`
}

type Seat struct {
	ID      uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	Name    string
	Status  SeatStatus
	EventID uuid.UUID `gorm:"not null"`
	Event   Event
}

type Reservation struct {
	ID      uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	UserID  uuid.UUID
	EventID uuid.UUID `gorm:"not null"`
	SeatID  uuid.UUID `gorm:"not null"`
	common.CreatedDeleted

	Event Event `gorm:"foreignKey:event_id"`
	Seat  Seat  `gorm:"foreignKey:seat_id"`
}

type EventRepository interface {
	CreateEvent(ctx context.Context, event *Event) error
	GetEventWithSeats(ctx context.Context, id uuid.UUID) (*Event, error)
	Migrate() error
}

type SeatRepository interface {
	CreateSeats(ctx context.Context, seats []*Seat) error
	GetSeat(ctx context.Context, id uuid.UUID) (*Seat, error)
	UpdateSeat(ctx context.Context, seat *Seat) error
	Migrate() error
}

type ReservationRepository interface {
	CreateReservation(ctx context.Context, reservation *Reservation) error
	GetReservationsForEvent(ctx context.Context, eventID uuid.UUID) ([]Reservation, error)
	GetUserReservations(ctx context.Context, userID uuid.UUID) ([]Reservation, error)
	GetReservationForSeat(ctx context.Context, seatID uuid.UUID) (*Reservation, error)
	GetReservationForEventAndUser(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) (*Reservation, error)
	HandleWithTransaction(ctx context.Context, fn ReservationTransaction) error
	DeleteReservation(ctx context.Context, id uuid.UUID) error
	Migrate() error
}

type SeatStatus string

const (
	SeatStatusAvailable SeatStatus = "available"
	SeatStatusReserved  SeatStatus = "reserved"
)

type ReservationTransaction func(ctx context.Context, reservationRepo ReservationRepository) (bool, error)

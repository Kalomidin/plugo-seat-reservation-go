package event_manager

import (
	"context"
	"seat-reservation/pkg/manager/event/repository"

	"github.com/google/uuid"
)

type CreateEventRequest struct {
	Name      string
	SeatCount int
	CreatorID uuid.UUID
}

type Seat struct {
	ID     uuid.UUID
	Name   string
	Status repository.SeatStatus
}

type GetEventRequest struct {
	EventID uuid.UUID
}

type Event struct {
	ID        uuid.UUID
	Name      string
	CreatorID uuid.UUID
	Seats     []Seat
}

type CreateReservationRequest struct {
	EventID uuid.UUID
	SeatID  uuid.UUID
	UserID  uuid.UUID
}

type CreateReservationResponse struct {
	ReservationID uuid.UUID
	EventID       uuid.UUID
	Seat          Seat
}

type ConfirmReservationRequest struct {
	EventID uuid.UUID
	UserID  uuid.UUID
	SeatID  uuid.UUID
}

type ConfirmReservationResponse struct {
	EventID    uuid.UUID
	UserID     uuid.UUID
	SeatID     uuid.UUID
	IsReserved bool
}

type CancelReservationRequest struct {
	EventID uuid.UUID
	UserID  uuid.UUID
}

type CancelReservationResponse struct {
	EventID uuid.UUID
	UserID  uuid.UUID
}

type EventManager interface {
	CreateEvent(ctx context.Context, req CreateEventRequest) (*Event, error)
	GetEvent(ctx context.Context, req GetEventRequest) (*Event, error)
	CreateReservation(ctx context.Context, req CreateReservationRequest) (*CreateReservationResponse, error)
	ConfirmReservation(ctx context.Context, req ConfirmReservationRequest) (*ConfirmReservationResponse, error)
	CancelReservation(ctx context.Context, req CancelReservationRequest) (*CancelReservationResponse, error)
}

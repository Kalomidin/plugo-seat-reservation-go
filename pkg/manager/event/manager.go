package event_manager

import (
	"context"
	"fmt"
	"seat-reservation/pkg/manager/event/repository"

	"gorm.io/gorm"
)

type eventManager struct {
	eventRepo       repository.EventRepository
	seatRepo        repository.SeatRepository
	reservationRepo repository.ReservationRepository
}

func NewEventManager(eventRepo repository.EventRepository, seatRepo repository.SeatRepository, reservationRepo repository.ReservationRepository) EventManager {
	return &eventManager{
		eventRepo,
		seatRepo,
		reservationRepo,
	}
}

func (m *eventManager) CreateEvent(ctx context.Context, req CreateEventRequest) (*Event, error) {
	newEvent := repository.Event{
		Name:      req.Name,
		CreatorID: req.CreatorID,
	}
	if err := m.eventRepo.CreateEvent(ctx, &newEvent); err != nil {
		return nil, err
	}
	var repoSeats []*repository.Seat = make([]*repository.Seat, req.SeatCount)
	for i := 0; i < req.SeatCount; i++ {
		repoSeats[i] = &repository.Seat{
			EventID: newEvent.ID,
			Name:    fmt.Sprintf("Seat %+v", i),
			Status:  repository.SeatStatusAvailable,
		}
	}
	if err := m.seatRepo.CreateSeats(ctx, repoSeats); err != nil {
		return nil, err
	}

	var seats []Seat = make([]Seat, len(repoSeats))
	for i, seat := range repoSeats {
		seats[i] = Seat{
			ID:     seat.ID,
			Name:   seat.Name,
			Status: seat.Status,
		}
	}

	return &Event{
		ID:        newEvent.ID,
		Name:      newEvent.Name,
		CreatorID: newEvent.CreatorID,
		Seats:     seats,
	}, nil
}

func (m *eventManager) GetEvent(ctx context.Context, req GetEventRequest) (*Event, error) {
	event, err := m.eventRepo.GetEventWithSeats(ctx, req.EventID)
	if err != nil {
		return nil, err
	}

	var seats []Seat = make([]Seat, len(event.Seats))
	for i, seat := range event.Seats {
		seats[i] = Seat{
			ID:   seat.ID,
			Name: seat.Name,
		}
	}
	return &Event{
		ID:        event.ID,
		Name:      event.Name,
		CreatorID: event.CreatorID,
		Seats:     seats,
	}, nil
}

func (m *eventManager) CreateReservation(ctx context.Context, req CreateReservationRequest) (*CreateReservationResponse, error) {
	var res repository.Reservation
	var seat repository.Seat
	if err := m.reservationRepo.HandleWithTransaction(ctx, func(ctx context.Context, db *gorm.DB) error {
		existingReservation, err := m.reservationRepo.GetReservationForEventAndUser(ctx, req.EventID, req.UserID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if existingReservation != nil || err != gorm.ErrRecordNotFound {
			return fmt.Errorf("user already has reservation for this event")
		}

		_seat, err := m.seatRepo.GetSeat(ctx, req.SeatID)
		if err != nil {
			return err
		}
		if _seat.Status != repository.SeatStatusAvailable {
			return fmt.Errorf("seat is not available")
		}
		_seat.Status = repository.SeatStatusReserved

		if err := m.seatRepo.UpdateSeat(ctx, _seat); err != nil {
			return err
		}
		seat = *_seat

		reservation := repository.Reservation{
			EventID: req.EventID,
			SeatID:  req.SeatID,
			UserID:  req.UserID,
		}
		if err := m.reservationRepo.CreateReservation(ctx, &reservation); err != nil {
			return err
		}
		res = reservation
		return nil
	}); err != nil {
		return nil, err
	}
	return &CreateReservationResponse{
		ReservationID: res.ID,
		EventID:       res.EventID,
		Seat: Seat{
			ID:     seat.ID,
			Name:   seat.Name,
			Status: seat.Status,
		},
	}, nil
}

func (m *eventManager) ConfirmReservation(ctx context.Context, req ConfirmReservationRequest) (*ConfirmReservationResponse, error) {
	reservation, err := m.reservationRepo.GetReservationForSeat(ctx, req.SeatID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &ConfirmReservationResponse{
				EventID:    req.EventID,
				UserID:     req.UserID,
				SeatID:     req.SeatID,
				IsReserved: false,
			}, nil
		}

		return nil, err
	}
	if reservation.UserID != req.UserID {
		return nil, fmt.Errorf("user is not owner of reservation")
	}
	if reservation.EventID != req.EventID {
		return nil, fmt.Errorf("reservation is not for this event")
	}
	return &ConfirmReservationResponse{
		EventID:    reservation.EventID,
		UserID:     reservation.UserID,
		SeatID:     reservation.SeatID,
		IsReserved: true,
	}, nil
}

func (m *eventManager) CancelReservation(ctx context.Context, req CancelReservationRequest) (*CancelReservationResponse, error) {
	if err := m.reservationRepo.HandleWithTransaction(ctx, func(ctx context.Context, db *gorm.DB) error {
		reservation, err := m.reservationRepo.GetReservationForEventAndUser(ctx, req.EventID, req.UserID)
		if err != nil {
			return err
		}
		if err := m.reservationRepo.DeleteReservation(ctx, reservation.ID); err != nil {
			return err
		}
		seat, err := m.seatRepo.GetSeat(ctx, reservation.SeatID)
		if err != nil {
			return err
		}
		seat.Status = repository.SeatStatusAvailable
		if err := m.seatRepo.UpdateSeat(ctx, seat); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &CancelReservationResponse{
		EventID: req.EventID,
		UserID:  req.UserID,
	}, nil
}

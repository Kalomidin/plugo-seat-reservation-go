package event_handler

import (
	"seat-reservation/common"
	event_manager "seat-reservation/pkg/manager/event"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateReservationRequest struct {
	SeatID uuid.UUID `json:"seatId"`
}

type CreateReservationResponse struct {
	ReservationID uuid.UUID `json:"reservationId"`
	EventID       uuid.UUID `json:"eventId"`
	Seat
}

func (h *handler) CreateReservation(ctx *gin.Context, req CreateReservationRequest) (*CreateReservationResponse, error) {
	eventID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return nil, err
	}
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, nil
	}

	manReq := event_manager.CreateReservationRequest{
		EventID: eventID,
		SeatID:  req.SeatID,
		UserID:  userID,
	}
	res, err := h.manager.CreateReservation(ctx, manReq)
	if err != nil {
		return nil, err
	}
	return &CreateReservationResponse{
		ReservationID: res.ReservationID,
		EventID:       res.EventID,
		Seat: Seat{
			ID:     res.Seat.ID,
			Name:   res.Seat.Name,
			Status: string(res.Seat.Status),
		},
	}, nil
}

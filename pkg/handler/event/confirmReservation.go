package event_handler

import (
	"net/http"
	"seat-reservation/common"
	event_manager "seat-reservation/pkg/manager/event"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ConfirmReservationResponse struct {
	EventID    uuid.UUID `json:"event_id"`
	UserID     uuid.UUID `json:"user_id"`
	SeatID     uuid.UUID `json:"seat_id"`
	IsReserved bool      `json:"is_reserved"`
}

func (h *HttpHandler) ConfirmReservation(ctx *gin.Context, req *http.Request) (interface{}, error) {
	resp, err := h.handler.ConfirmReservation(ctx)
	return resp, err
}

func (h *handler) ConfirmReservation(ctx *gin.Context) (*ConfirmReservationResponse, error) {
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	eventID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return nil, err
	}
	seatID, err := uuid.Parse(ctx.Param("seatId"))
	if err != nil {
		return nil, err
	}
	req := event_manager.ConfirmReservationRequest{
		EventID: eventID,
		UserID:  userID,
		SeatID:  seatID,
	}

	res, err := h.manager.ConfirmReservation(ctx, req)
	if err != nil {
		return nil, err
	}
	return &ConfirmReservationResponse{
		EventID:    res.EventID,
		UserID:     res.UserID,
		SeatID:     res.SeatID,
		IsReserved: res.IsReserved,
	}, nil
}

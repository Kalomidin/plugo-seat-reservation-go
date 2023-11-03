package event_handler

import (
	"net/http"
	"seat-reservation/common"
	event_manager "seat-reservation/pkg/manager/event"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CancelReservationResponse struct {
	EventID uuid.UUID `json:"event_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func (h *HttpHandler) CancelReservation(ctx *gin.Context, req *http.Request) (interface{}, error) {
	resp, err := h.handler.CancelReservation(ctx)
	return resp, err
}

func (h *handler) CancelReservation(ctx *gin.Context) (*CancelReservationResponse, error) {
	eventId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return nil, err
	}
	userId, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	req := event_manager.CancelReservationRequest{
		EventID: eventId,
		UserID:  userId,
	}
	resp, err := h.manager.CancelReservation(ctx, req)
	if err != nil {
		return nil, err
	}
	return &CancelReservationResponse{
		EventID: resp.EventID,
		UserID:  resp.UserID,
	}, nil
}

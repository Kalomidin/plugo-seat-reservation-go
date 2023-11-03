package event_handler

import (
	"net/http"
	"seat-reservation/common"
	event_manager "seat-reservation/pkg/manager/event"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateEventRequest struct {
	Name      string `json:"name"`
	SeatCount int    `json:"seatCount"`
}

type CreateEventResponse struct {
	Event
}

type Event struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatorID uuid.UUID `json:"creatorId"`
	Seats     []Seat    `json:"seats"`
}

type Seat struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Status string    `json:"status"`
}

func (h *HttpHandler) CreateEvent(ctx *gin.Context, req *http.Request) (interface{}, error) {
	var createEventRequest CreateEventRequest
	if err := ctx.ShouldBindJSON(&createEventRequest); err != nil {
		return nil, err
	}
	resp, err := h.handler.CreateEvent(ctx, createEventRequest)
	return resp, err
}

func (h *handler) CreateEvent(ctx *gin.Context, req CreateEventRequest) (*CreateEventResponse, error) {
	userId, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}

	manReq := event_manager.CreateEventRequest{
		Name:      req.Name,
		SeatCount: req.SeatCount,
		CreatorID: userId,
	}
	manResp, err := h.manager.CreateEvent(ctx, manReq)
	if err != nil {
		return nil, err
	}

	var seats []Seat = make([]Seat, len(manResp.Seats))
	for i, seat := range manResp.Seats {
		seats[i] = Seat{
			ID:     seat.ID,
			Name:   seat.Name,
			Status: string(seat.Status),
		}
	}

	return &CreateEventResponse{
		Event: Event{
			ID:        manResp.ID,
			Name:      manResp.Name,
			CreatorID: manResp.CreatorID,
			Seats:     seats,
		},
	}, nil
}

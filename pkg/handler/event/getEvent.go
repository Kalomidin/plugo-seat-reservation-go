package event_handler

import (
	"net/http"
	event_manager "seat-reservation/pkg/manager/event"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetEventResponse struct {
	Event
}

func (h *HttpHandler) GetEvent(ctx *gin.Context, req *http.Request) (interface{}, error) {
	resp, err := h.handler.GetEvent(ctx)
	return resp, err
}

func (h *handler) GetEvent(ctx *gin.Context) (*GetEventResponse, error) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return nil, err
	}

	manResp, err := h.manager.GetEvent(ctx, event_manager.GetEventRequest{
		EventID: id,
	})
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

	return &GetEventResponse{
		Event: Event{
			ID:        manResp.ID,
			Name:      manResp.Name,
			CreatorID: manResp.CreatorID,
			Seats:     seats,
		},
	}, nil

}

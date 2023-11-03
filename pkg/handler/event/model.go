package event_handler

import "github.com/gin-gonic/gin"

type EventHandler interface {
	CreateEvent(ctx *gin.Context, req CreateEventRequest) (*CreateEventResponse, error)
	GetEvent(ctx *gin.Context) (*GetEventResponse, error)
	CreateReservation(ctx *gin.Context, req CreateReservationRequest) (*CreateReservationResponse, error)
	ConfirmReservation(ctx *gin.Context) (*ConfirmReservationResponse, error)
	CancelReservation(ctx *gin.Context) (*CancelReservationResponse, error)
}

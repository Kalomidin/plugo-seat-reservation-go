package event_handler

import (
	"context"
	"fmt"
	"seat-reservation/common"
	manager "seat-reservation/pkg/manager/event"

	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	handler    EventHandler
	middleware common.Middleware
}

type handler struct {
	manager manager.EventManager
}

func NewHandler(manager manager.EventManager) EventHandler {
	return &handler{
		manager,
	}
}

func NewHttpHandler(ctx context.Context, h EventHandler, middleware common.Middleware) *HttpHandler {
	return &HttpHandler{
		h,
		middleware,
	}
}

func (c *HttpHandler) Init(ctx context.Context, router *gin.Engine) {
	routes := map[string]map[string]common.HandlerFunc{
		"POST": {
			"": c.middleware.HandlerWithAuth(c.CreateEvent),
		},
		"GET": {
			"/:id": c.middleware.HandlerWithAuth(c.GetEvent),
			"/:id/reservation/confirm?seatId=:seatId": c.middleware.HandlerWithAuth(c.ConfirmReservation),
		},
		"DELETE": {
			"/:id": c.middleware.HandlerWithAuth(c.CancelReservation),
		},
	}
	for method, route := range routes {
		for r, h := range route {
			router.Handle(method, fmt.Sprintf("/event%s", r), common.GenericHandler(h))
		}
	}
	fmt.Println("initialized event handler")
}

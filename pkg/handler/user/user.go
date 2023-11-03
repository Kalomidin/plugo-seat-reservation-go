package user_handler

import (
	"context"
	"fmt"
	"seat-reservation/common"
	manager "seat-reservation/pkg/manager/user"

	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	handler    UserHandler
	middleware common.Middleware
}

type handler struct {
	manager manager.UserManager
}

func NewHandler(manager manager.UserManager) UserHandler {
	return &handler{
		manager,
	}
}

func NewHttpHandler(ctx context.Context, h UserHandler, middleware common.Middleware) *HttpHandler {
	return &HttpHandler{
		h,
		middleware,
	}
}

func (c *HttpHandler) Init(ctx context.Context, router *gin.Engine) {
	routes := map[string]map[string]common.HandlerFunc{
		"GET": {
			"": c.middleware.HandlerWithAuth(c.GetUser),
		},
	}
	for method, route := range routes {
		for r, h := range route {
			router.Handle(method, fmt.Sprintf("/user%s", r), common.GenericHandler(h))
		}
	}
	fmt.Println("initialized user handler")
}

package user_handler

import (
	"net/http"
	"seat-reservation/common"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetUserResponse struct {
	UserId   uuid.UUID `json:"id"`
	Username string    `json:"userName"`
}

func (h *HttpHandler) GetUser(ctx *gin.Context, r *http.Request) (interface{}, error) {
	resp, err := h.handler.GetUser(ctx)
	return resp, err
}

func (h *handler) GetUser(ctx *gin.Context) (*GetUserResponse, error) {
	userID, err := common.GetUserId(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	user, err := h.manager.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &GetUserResponse{
		UserId:   user.ID,
		Username: user.Username,
	}, nil
}

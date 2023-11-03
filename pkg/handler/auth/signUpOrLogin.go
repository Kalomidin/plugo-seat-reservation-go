package auth_handler

import (
	"net/http"
	auth_manager "seat-reservation/pkg/manager/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SignupOrLoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type SignupOrLoginResponse struct {
	Id           uuid.UUID `json:"id"`
	UserName     string    `json:"userName"`
	AuthToken    string    `json:"authToken"`
	RefreshToken string    `json:"refreshToken"`
}

func (h *HttpHandler) SignupOrLogin(ctx *gin.Context, r *http.Request) (interface{}, error) {
	var req SignupOrLoginRequest
	if err := ctx.BindJSON(&req); err != nil {
		return nil, err
	}
	resp, err := h.handler.SignupOrLogin(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *handler) SignupOrLogin(ctx *gin.Context, req SignupOrLoginRequest) (*SignupOrLoginResponse, error) {
	manReq := auth_manager.SignupOrLoginRequest{
		UserName: req.UserName,
		Password: req.Password,
	}

	manResp, err := h.service.SignupOrLogin(ctx, manReq)
	if err != nil {
		return nil, err
	}
	return &SignupOrLoginResponse{
		Id:           manResp.Id,
		UserName:     manResp.UserName,
		AuthToken:    manResp.AuthToken,
		RefreshToken: manResp.RefreshToken,
	}, nil
}

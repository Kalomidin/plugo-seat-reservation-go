package auth_handler

import (
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	SignupOrLogin(ctx *gin.Context, req SignupOrLoginRequest) (*SignupOrLoginResponse, error)
}

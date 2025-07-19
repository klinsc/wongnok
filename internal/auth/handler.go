package auth

import "github.com/gin-gonic/gin"

type IHandler interface {
}

type Handler struct {
	Service IService
}

func NewHandler() IHandler {
	return &Handler{Service: NewService()}
}

func (handler Handler) Login(ctx *gin.Context) {
	// Generate a state parameter for CSRF protection
	state := handler.Service.GenerateState()

	// Collect state in Cookie
	ctx.SetCookie("state", state, 300, "/", "localhost", false, true)
}

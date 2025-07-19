package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IHandler interface {
	Login(ctx *gin.Context)
}

type Handler struct {
	Service IService
}

func NewHandler(oauth2Conf IOAuth2Config) IHandler {
	return &Handler{
		Service: NewService(oauth2Conf),
	}
}

func (handler Handler) Login(ctx *gin.Context) {
	// Generate state
	state := handler.Service.GenerateState()

	// Collect state in Cookie
	ctx.SetCookie("state", state, 300, "/", "localhost", false, true)

	fmt.Println(handler.Service.AuthCodeURL(state))

	// Redirect to Keycloak login page
	ctx.JSON(http.StatusOK, gin.H{"message": true, "url": handler.Service.AuthCodeURL(state)})

}

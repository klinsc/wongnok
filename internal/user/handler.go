package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type IHandler interface{}

type Handler struct{}

func NewHandler() IHandler {
	return &Handler{}
}

func (handler Handler) GetRecipes(ctx *gin.Context) {
	tokenWithBearer := ctx.GetHeader("Authorization")
	fmt.Println("Token with Bearer:", tokenWithBearer)
}

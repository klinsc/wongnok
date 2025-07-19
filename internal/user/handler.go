package user

import (
	"github.com/gin-gonic/gin"
	"github.com/klins/devpool/go-day6/wongnok/internal/helper"
)

type IHandler interface {
	GetRecipes(ctx *gin.Context)
}

type Handler struct {
}

func NewHandler() IHandler {
	return &Handler{}
}

func (handler Handler) GetRecipes(ctx *gin.Context) {
	// tokenWithBearer := ctx.GetHeader("Authorization")
	// fmt.Println("Token with Bearer:", tokenWithBearer

	claims, err := helper.DecodeClaims(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "Invalid token"})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "User recipes",
		"user_id": claims.ID,
	})
}

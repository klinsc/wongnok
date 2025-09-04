package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/klins/devpool/go-day6/wongnok/internal/helper"
	"gorm.io/gorm"
)

type IHandler interface {
	GetRecipes(ctx *gin.Context)
	GetByID(ctx *gin.Context)
}

type Handler struct {
	Service IService
}

func NewHandler(db *gorm.DB) IHandler {
	return &Handler{
		Service: NewService(db),
	}
}

func (handler Handler) GetRecipes(ctx *gin.Context) {
	userID := ctx.Param("id")

	claims, err := helper.DecodeClaims(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	recipes, err := handler.Service.GetRecipes(userID, claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipes.ToResponse(int64(len(recipes))))
}

func (handler Handler) GetByID(ctx *gin.Context) {
	userID := ctx.Param("id")

	user, err := handler.Service.GetByID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user.ToResponse())
}


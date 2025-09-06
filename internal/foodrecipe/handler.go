package foodrecipe

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/klins/devpool/go-day6/wongnok/internal/global"
	"github.com/klins/devpool/go-day6/wongnok/internal/helper"
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"github.com/klins/devpool/go-day6/wongnok/internal/model/dto"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IHandler interface {
	Create(ctx *gin.Context)
	Get(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetFavorites(ctx *gin.Context)
}

type Handler struct {
	Service IService
}

func NewHandler(db *gorm.DB) IHandler {
	return &Handler{
		Service: NewService(db),
	}
}

func (handler Handler) Create(ctx *gin.Context) {
	var request dto.FoodRecipeRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	claims, err := helper.DecodeClaims(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	recipe, err := handler.Service.Create(request, claims)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.As(err, &validator.ValidationErrors{}) {
			statusCode = http.StatusBadRequest
		}

		ctx.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, recipe.ToResponse())
}

func (handler Handler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "ID is required"})
		return
	}

	recipe, err := handler.Service.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Recipe not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipe.ToResponse())
}

func (handler Handler) GetAll(ctx *gin.Context) {
	recipes, err := handler.Service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var response []dto.FoodRecipeResponse
	for _, recipe := range recipes {
		response = append(response, recipe.ToResponse())
	}

	ctx.JSON(http.StatusOK, response)
}

func (handler Handler) Get(ctx *gin.Context) {
	var foodRecipeQuery model.FoodRecipeQuery
	if err := ctx.ShouldBindQuery(&foodRecipeQuery); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	recipes, total, err := handler.Service.Get(foodRecipeQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipes.ToResponse(total))
}

func (handler Handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "ID is required"})
		return
	}

	var request dto.FoodRecipeRequest
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	claims, err := helper.DecodeClaims(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	recipe, err := handler.Service.Update(request, id, claims)
	if err != nil {
		if errors.Is(err, global.ErrorNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Recipe not found"})
			return
		}

		if errors.Is(err, global.ErrorForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "You do not have permission to update this recipe"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipe.ToResponse())
}

func (handler Handler) Delete(ctx *gin.Context) {
	claims, err := helper.DecodeClaims(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	var id string

	pathParam := ctx.Param("id")
	if pathParam != "" {
		if parsed, err := strconv.Atoi(pathParam); err == nil && parsed > 0 {
			id = strconv.Itoa(parsed)
		}
	}

	if err := handler.Service.Delete(id, claims); err != nil {
		statusCode := http.StatusInternalServerError

		if errors.Is(err, global.ErrorForbidden) {
			statusCode = http.StatusForbidden
		}

		ctx.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
}

func (handler Handler) GetFavorites(ctx *gin.Context) {
	claims, err := helper.DecodeClaims(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	var foodRecipeQuery model.FoodRecipeQuery
	if err := ctx.ShouldBindQuery(&foodRecipeQuery); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	recipes, total, err := handler.Service.GetFavorites(foodRecipeQuery, claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipes.ToResponse(total))
}

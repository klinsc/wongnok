package helper

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
)

func DecodeClaims(ctx *gin.Context) (model.Claims, error) {
	value, exists := ctx.Get("claims")
	if !exists {
		return model.Claims{}, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	claims, ok := value.(model.Claims)
	if !ok {
		return model.Claims{}, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	return claims, nil
}

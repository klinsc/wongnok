package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/klins/devpool/go-day6/wongnok/config"
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
)

func Authorize(
	verifier config.IOIDCTokenVerifier,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerPrefix := "Bearer "

		tokenWithBearer := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(tokenWithBearer, bearerPrefix) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing Bearer token"})
			ctx.Abort()
			return
		}

		rawToken := strings.TrimPrefix(tokenWithBearer, bearerPrefix)
		idToken, err := verifier.Verify(ctx, rawToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "error": err.Error()})
			ctx.Abort()
			return
		}

		var claims model.Claims
		if err := idToken.Claims(&claims); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Failed to parse token claims", "error": err.Error()})
			ctx.Abort()
			return
		}

		// Set user claims in context
		ctx.Set("claims", claims)

		// Continue to the next handler
		ctx.Next()
	}
}

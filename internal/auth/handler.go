package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/klins/devpool/go-day6/wongnok/config"
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"github.com/klins/devpool/go-day6/wongnok/internal/model/dto"
	"gorm.io/gorm"
)

type IHandler interface {
	Login(ctx *gin.Context)
	Callback(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type Handler struct {
	Service IService
}

func NewHandler(db *gorm.DB, kc config.Keycloak, oauth2Conf IOAuth2Config, verifier IOIDCTokenVerifier) IHandler {
	return &Handler{
		Service: NewService(kc, oauth2Conf, verifier),
	}
}

func (handler Handler) Login(ctx *gin.Context) {
	// Generate state
	state := handler.Service.GenerateState()

	// Collect state in Cookie
	ctx.SetCookie("state", state, 300, "/", "localhost", false, true)

	fmt.Println(handler.Service.AuthCodeURL(state))

	// Redirect to Keycloak
	ctx.Redirect(http.StatusTemporaryRedirect, handler.Service.AuthCodeURL(state))

}

func (handler Handler) Callback(ctx *gin.Context) {
	// Get state from Cookie
	_, err := ctx.Cookie("state")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "State cookie not found", "error": err.Error()})
		return
	}

	// Get code from query parameters
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Code not found in query parameters"})
		return
	}

	// Exchange code for token
	credential, err := handler.Service.Exchange(ctx.Request.Context(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to exchange code for token", "error": err.Error()})
		return
	}

	// Verify token
	idToken, err := handler.Service.VerifyToken(ctx.Request.Context(), credential.IDToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify token", "error": err.Error()})
		return
	}

	var claims model.Claims
	if err := idToken.Claims(&claims); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse token claims", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, credential.ToResponse())
}

func (handler Handler) Logout(ctx *gin.Context) {
	var query dto.LogoutQuery
	ctx.BindQuery(&query)

	// https://{domain}/realms/{realm}/protocol/openid-connect/logout?id_token_hint{idTokenHint}&post_logout_redirect_uri={postLogoutRedirectUri}

	// Make lougout URL
	logoutURL, err := handler.Service.LogoutURL(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create logout URL", "error": err.Error()})
		return
	}

	// Redirect to Keycloak logout page
	ctx.Redirect(http.StatusFound, logoutURL)

}

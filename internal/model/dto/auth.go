package dto

import "time"

type KeycloakCallbackQuery struct {
	State string `form:"state"`
	Code  string `form:"code"`
}

type CredentialResponse struct {
	AccessToken  string    `json:"accessToken"`
	TokenType    string    `json:"tokenType,omitempty"`
	RefreshToken string    `json:"refreshToken,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
	ExpiresIn    int64     `json:"expiresIn"`
	IDToken      string    `json:"idToken"`
}

type LogoutQuery struct {
	IDTokenHint           string `form:"idTokenHint"`
	PostLogoutRedirectURI string `form:"postLogoutRedirectUri"`
}

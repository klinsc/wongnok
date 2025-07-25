package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/coreos/go-oidc"
	"github.com/klins/devpool/go-day6/wongnok/config"
	"github.com/klins/devpool/go-day6/wongnok/internal/model"
	"github.com/klins/devpool/go-day6/wongnok/internal/model/dto"
	"github.com/pkg/errors"
)

type IOAuth2Config config.IOAuth2Config
type IOIDCTokenVerifier config.IOIDCTokenVerifier
type IOIDCIDToken config.IOIDCIDToken

type IService interface {
	GenerateState() string
	AuthCodeURL(state string) string
	Exchange(ctx context.Context, code string) (model.Credential, error)
	VerifyToken(ctx context.Context, tokenString string) (*oidc.IDToken, error)
	LogoutURL(logoutQuery dto.LogoutQuery) (string, error)
}

type Service struct {
	KeyCloak     config.Keycloak
	OAuth2Config IOAuth2Config
	Verifier     IOIDCTokenVerifier
}

func NewService(kc config.Keycloak, oauth2Config IOAuth2Config, verifier IOIDCTokenVerifier) IService {
	return &Service{
		KeyCloak:     kc,
		OAuth2Config: oauth2Config,
		Verifier:     verifier,
	}
}

func (service Service) GenerateState() string {
	buffer := make([]byte, 32)
	rand.Read(buffer)
	return base64.URLEncoding.EncodeToString(buffer)
}

func (service Service) AuthCodeURL(state string) string {
	return service.OAuth2Config.AuthCodeURL(state)
}

func (service Service) Exchange(ctx context.Context, code string) (model.Credential, error) {
	token, err := service.OAuth2Config.Exchange(ctx, code)
	if err != nil {
		return model.Credential{}, errors.Wrap(err, "exchange token")
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return model.Credential{}, fmt.Errorf("id token is missing")
	}

	return model.Credential{
		Token:   token,
		IDToken: idToken,
	}, nil
}

func (service Service) VerifyToken(ctx context.Context, token string) (*oidc.IDToken, error) {
	idToken, err := service.Verifier.Verify(ctx, token)
	if err != nil {
		return nil, errors.Wrap(err, "verify token")
	}

	return idToken, nil
}

// https://{domain}/realms/{realm}/protocol/openid-connect/logout?id_token_hint{idTokenHint}&post_logout_redirect_uri={postLogoutRedirectUri}
func (service Service) LogoutURL(logoutQuery dto.LogoutQuery) (string, error) {
	uri, err := url.Parse(service.KeyCloak.LogoutURL())
	if err != nil {
		return "", errors.Wrap(err, "parse logout URL")
	}

	query := uri.Query()
	query.Set("id_token_hint", logoutQuery.IDTokenHint)
	query.Set("post_logout_redirect_uri", logoutQuery.PostLogoutRedirectURI)
	uri.RawQuery = query.Encode()

	return uri.String(), nil
}

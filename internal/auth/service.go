package auth

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/klins/devpool/go-day6/wongnok/config"
)

type IOAuth2Config config.IOAuth2Config

type IService interface {
	GenerateState() string
	AuthCodeURL(state string) string
}

type Service struct {
	OAuth2Config IOAuth2Config
}

func NewService(oauth2Config IOAuth2Config) IService {
	return &Service{
		OAuth2Config: oauth2Config,
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

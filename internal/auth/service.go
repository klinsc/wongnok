package auth

import (
	"crypto/rand"
	"encoding/base64"
)

type IService interface {
	GenerateState() string
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GenerateState() string {
	buffer := make([]byte, 32)
	rand.Read(buffer)
	return base64.URLEncoding.EncodeToString(buffer)
}

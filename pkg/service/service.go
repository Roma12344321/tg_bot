package service

import (
	"context"
	"myBot/pkg/storage"
)

const BaseURL = "http://romamar2004.fvds.ru/api"

type AuthService interface {
	LogIn(telegramId int64, username, password string) error
}

type Service struct {
	AuthService
}

func NewService(storage *storage.Storage, ctx context.Context) *Service {
	return &Service{
		AuthService: NewAuthServiceImpl(storage, ctx),
	}
}

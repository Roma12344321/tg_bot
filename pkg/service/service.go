package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"myBot/pkg/model"
	"myBot/pkg/storage"
	"net/http"
	"time"
)

const (
	baseURL = "http://romamar2004.fvds.ru/api"
	bearer  = "Bearer "
)

type AuthService interface {
	LogIn(telegramId int64, username, password string) error
}

type GoodService interface {
	GetAllGood(telegramId int64) ([]model.Good, error)
}

type CartService interface {
	GetAllCart(telegramId int64) ([]model.Cart, error)
}

type Service struct {
	AuthService
	GoodService
	CartService
}

func NewService(ctx context.Context, storage *storage.Storage, client *http.Client) *Service {
	return &Service{
		AuthService: NewAuthServiceImpl(ctx, storage),
		GoodService: NewGoodServiceImpl(ctx, storage, client),
		CartService: NewCartServiceImpl(ctx, storage, client),
	}
}

func handleGetRequestWithTokenReturningList[T any](
	ctx context.Context,
	storage *storage.Storage,
	telegramId int64,
	client *http.Client,
	url string) ([]T, error) {

	tokenFromMap, ok := storage.Tokens.Load(telegramId)
	if !ok {
		return nil, errors.New("not authorized")
	}
	token, ok := tokenFromMap.(string)
	if !ok {
		log.Fatalln("map error")
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", bearer+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	objects := make([]T, 0)
	err = json.NewDecoder(resp.Body).Decode(&objects)
	if err != nil {
		return nil, err
	}
	return objects, nil
}

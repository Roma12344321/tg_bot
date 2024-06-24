package service

import (
	"context"
	"myBot/pkg/model"
	"myBot/pkg/storage"
	"net/http"
)

type CartServiceImpl struct {
	ctx     context.Context
	storage *storage.Storage
	client  *http.Client
}

func NewCartServiceImpl(ctx context.Context, storage *storage.Storage, client *http.Client) *CartServiceImpl {
	return &CartServiceImpl{ctx: ctx, storage: storage, client: client}
}

func (s *CartServiceImpl) GetAllCart(telegramId int64) ([]model.Cart, error) {
	return handleGetRequestWithTokenReturningList[model.Cart](s.ctx, s.storage, telegramId, s.client, "/cart")
}

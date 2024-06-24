package service

import (
	"context"
	"myBot/pkg/model"
	"myBot/pkg/storage"
	"net/http"
)

type GoodServiceImpl struct {
	ctx     context.Context
	storage *storage.Storage
	client  *http.Client
}

func NewGoodServiceImpl(ctx context.Context, storage *storage.Storage, client *http.Client) *GoodServiceImpl {
	return &GoodServiceImpl{ctx: ctx, storage: storage, client: client}
}

func (s *GoodServiceImpl) GetAllGood(telegramId int64) ([]model.Good, error) {
	return handleGetRequestWithTokenReturningList[model.Good](s.ctx, s.storage, telegramId, s.client, "/good")
}

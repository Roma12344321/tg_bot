package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"myBot/pkg/model"
	"myBot/pkg/storage"
	"net/http"
	"time"
)

type AuthServiceImpl struct {
	storage *storage.Storage
	ctx     context.Context
}

func NewAuthServiceImpl(storage *storage.Storage, ctx context.Context) *AuthServiceImpl {
	return &AuthServiceImpl{storage: storage, ctx: ctx}
}

func (s *AuthServiceImpl) LogIn(telegramId int64, username, password string) error {
	token, err := s.logIn(username, password)
	if err != nil {
		return err
	}
	s.storage.Tokens.Store(telegramId, token)
	return nil
}

func (s *AuthServiceImpl) logIn(username, password string) (string, error) {
	user := model.User{Username: username, Password: password}
	data, err := json.Marshal(&user)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()
	ch := make(chan struct {
		token string
		err   error
	}, 1)
	go func() {
		token, err := s.getToken(data)
		ch <- struct {
			token string
			err   error
		}{token, err}
		close(ch)
	}()
	select {
	case <-ctx.Done():
		return "", errors.New("request timeout")
	case result := <-ch:
		return result.token, result.err
	}
}

func (s *AuthServiceImpl) getToken(data []byte) (string, error) {
	resp, err := http.Post(BaseURL+"/auth/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	var token model.Token
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return "", err
	}
	if len(token.Token) < 30 {
		return "", errors.New("invalid username or password")
	}
	return token.Token, nil
}

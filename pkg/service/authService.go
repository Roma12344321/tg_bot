package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"myBot/pkg/model"
	"myBot/pkg/storage"
	"net/http"
	"time"
)

type AuthServiceImpl struct {
	ctx     context.Context
	storage *storage.Storage
}

func NewAuthServiceImpl(ctx context.Context, storage *storage.Storage) *AuthServiceImpl {
	return &AuthServiceImpl{ctx: ctx, storage: storage}
}

func (s *AuthServiceImpl) LogIn(telegramId int64, username, password string) error {
	token, err := s.logIn(username, password)
	if err != nil {
		return err
	}
	s.storage.Tokens.Store(telegramId, token)
	return nil
}

type tokenWithErr struct {
	token string
	err   error
}

func (s *AuthServiceImpl) logIn(username, password string) (string, error) {
	user := model.User{Username: username, Password: password}
	data, err := json.Marshal(&user)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()
	ch := make(chan tokenWithErr, 1)
	go func() {
		ch <- s.getToken(data)
		close(ch)
	}()
	select {
	case <-ctx.Done():
		log.Println("request was cancelled because of context...")
		return "", errors.New("request timeout")
	case result := <-ch:
		return result.token, result.err
	}
}

func (s *AuthServiceImpl) getToken(data []byte) tokenWithErr {
	resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return tokenWithErr{token: "", err: err}
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	var token model.Token
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return tokenWithErr{token: "", err: err}
	}
	if len(token.Token) < 30 {
		return tokenWithErr{token: "", err: errors.New("invalid username or password")}
	}
	return tokenWithErr{token: token.Token, err: nil}
}

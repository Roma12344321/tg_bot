package storage

import "sync"

type Storage struct {
	Tokens *sync.Map
}

func NewStorage(tokens *sync.Map) *Storage {
	return &Storage{Tokens: tokens}
}

package redis

import (
	"context"
	"time"
)

type Service struct {
	client *Client
}

func NewService(client *Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) Set(
	ctx context.Context,
	key string,
	value interface{},
	ttl time.Duration,
) error {
	return s.client.rdb.Set(
		ctx,
		key,
		value,
		ttl,
	).Err()
}

func (s *Service) Get(
	ctx context.Context,
	key string,
) (string, error) {
	return s.client.rdb.Get(
		ctx,
		key,
	).Result()
}

func (s *Service) Delete(
	ctx context.Context,
	keys ...string,
) error {
	return s.client.rdb.Del(
		ctx,
		keys...,
	).Err()
}

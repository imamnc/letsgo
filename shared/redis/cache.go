package redis

import (
	"context"
	"time"
)

// Remember Cache
func (s *Service) Remember(
	ctx context.Context,
	key string,
	ttl time.Duration,
	callback func() (string, error),
) (string, error) {

	val, err := s.Get(ctx, key)

	if err == nil {
		return val, nil
	}

	val, err = callback()

	if err != nil {
		return "", err
	}

	s.Set(ctx, key, val, ttl)

	return val, nil
}

// Forget Cache
func (s *Service) Forget(
	ctx context.Context,
	key string,
) error {
	return s.Delete(ctx, key)
}

// ClearAll Cache
func (s *Service) ClearAll(ctx context.Context) error {
	return s.client.rdb.FlushDB(ctx).Err()
}

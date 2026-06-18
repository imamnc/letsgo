package redis

import (
	"context"
	"time"
)

func (s *Service) AcquireLock(
	ctx context.Context,
	key string,
	ttl time.Duration,
) bool {

	ok, _ := s.client.rdb.SetNX(
		ctx,
		key,
		1,
		ttl,
	).Result()

	return ok
}

func (s *Service) ReleaseLock(
	ctx context.Context,
	key string,
) {
	s.client.rdb.Del(ctx, key)
}

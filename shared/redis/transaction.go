package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

func (s *Service) Transaction(
	ctx context.Context,
	fn func(pipe goredis.Pipeliner) error,
) error {

	_, err := s.client.rdb.TxPipelined(
		ctx,
		func(pipe goredis.Pipeliner) error {
			return fn(pipe)
		},
	)

	return err
}

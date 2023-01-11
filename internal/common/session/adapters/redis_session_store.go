package adapters

import (
	"context"

	"github.com/BON4/gofeed/internal/common/errors"
	"github.com/BON4/gofeed/internal/common/session/domain"
	"github.com/go-redis/redis/v9"
)

type RedisStore struct {
	client *redis.Client
	fc     *domain.SessionFactory
}

func NewRedisConnection(host string, password string, db int) (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	err := cli.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func NewRedisStore(client *redis.Client, fc *domain.SessionFactory) *RedisStore {
	return &RedisStore{
		client: client,
		fc:     fc,
	}
}

func (r *RedisStore) Set(ctx context.Context, key string, ss *domain.Session) error {
	marshaled, err := ss.MarshalJSON()
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, marshaled, ss.GetTTL()).Err()
}

func (r *RedisStore) Get(ctx context.Context, key string) (*domain.Session, error) {
	res, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, errors.NewDoesNotExistsError("session does not exists", "session-not-found")
	}

	ss, err := r.fc.UnmarshalSessionJSON([]byte(res))
	if err != nil {
		return nil, err
	}

	return ss, nil
}

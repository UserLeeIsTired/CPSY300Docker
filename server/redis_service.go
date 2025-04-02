package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Redis struct {
	client *redis.Client
}

func NewRedis() *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	})

	return &Redis{client: rdb}
}

func (r *Redis) Ping() error {
	pong, err := r.client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	fmt.Println("Ping result:", pong)
	return nil
}

func (r *Redis) CreateKeyWithExpiration() (string, error) {
	key := uuid.New().String()
	value := uuid.New().String()
	err := r.client.Set(context.Background(), key, value, 1*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return key, nil
}

func (r *Redis) GetValueByKey(key string) (string, error) {
	result, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}

	err = r.ExtendKeyExpiration(key)

	if err != nil {
		return "", err
	}

	return result, nil
}

func (r *Redis) DeleteKey(key string) (bool, error) {
	err := r.client.Del(context.Background(), key).Err()

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *Redis) ExtendKeyExpiration(key string) error {
	ttlResult, err := r.client.TTL(context.Background(), key).Result()
	if err != nil {
		return err
	}

	if ttlResult >= 0 {
		err := r.client.Expire(context.Background(), key, 1*time.Hour).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

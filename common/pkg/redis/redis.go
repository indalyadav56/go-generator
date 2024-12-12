package redis

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Publish(channel string, message string) *redis.IntCmd
	Subscribe(channel string) *redis.PubSub
	Unsubscribe(channel string, pubsub *redis.PubSub) error
	SAdd(key string, members ...interface{}) *redis.IntCmd
	SMembers(key string) *redis.StringSliceCmd
	SRem(key string, members ...interface{}) *redis.IntCmd
}

type redisService struct {
	client *redis.Client
}

type Config struct {
	Addr     string
	Username string
	Password string
	DB       int
	SSL      bool
}

func New(config Config) (*redisService, error) {
	options := &redis.Options{
		Addr:     config.Addr,
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	}

	if config.SSL {
		options.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
		}
	}

	client := redis.NewClient(options)

	// Test the connection using Ping
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()

	// Handle errors during connection
	if err != nil {
		log.Printf("Failed to connect to Redis: %v\n", err)
		return nil, err
	}

	// Return the Redis service if successful
	return &redisService{client: client}, nil
}

func (r *redisService) Get(key string) *redis.StringCmd {
	return r.client.Get(context.Background(), key)
}

func (r *redisService) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(context.Background(), key, value, expiration)
}

func (r *redisService) Del(keys ...string) *redis.IntCmd {
	return r.client.Del(context.Background(), keys...)
}

func (r *redisService) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	return r.client.Expire(context.Background(), key, expiration)
}

func (r *redisService) SAdd(key string, members ...interface{}) *redis.IntCmd {
	return r.client.SAdd(context.Background(), key, members)
}

func (r *redisService) SMembers(key string) *redis.StringSliceCmd {
	return r.client.SMembers(context.Background(), key)
}

func (r *redisService) SRem(key string, members ...interface{}) *redis.IntCmd {
	return r.client.SRem(context.Background(), key, members)
}

func (r *redisService) LPush(key string, values ...interface{}) *redis.IntCmd {
	return r.client.LPush(context.Background(), key, values...)
}

func (r *redisService) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	return r.client.LRange(context.Background(), key, start, stop)
}

func (r *redisService) LPop(key string) *redis.StringCmd {
	return r.client.LPop(context.Background(), key)
}

func (r *redisService) HSet(key, field, value string) {
	r.client.HSet(context.Background(), key, field, value)
}

func (r *redisService) HGetAll(key string) *redis.StringStringMapCmd {
	return r.client.HGetAll(context.Background(), key)
}

func (r *redisService) XAdd(stream string, values map[string]interface{}) {
	r.client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: stream,
		Values: values,
	})
}

func (r *redisService) XRange(stream, start, stop string) *redis.XMessageSliceCmd {
	return r.client.XRange(context.Background(), stream, start, stop)
}

func (r *redisService) SetJson(key string, value interface{}) {
	json, _ := json.Marshal(value)
	r.client.Set(context.Background(), key, json, 0)
}

func (r *redisService) Publish(channel string, message string) *redis.IntCmd {
	return r.client.Publish(context.Background(), channel, message)
}

func (r *redisService) Subscribe(channel string) *redis.PubSub {
	return r.client.Subscribe(context.Background(), channel)
}

func (r *redisService) Unsubscribe(channel string, pubsub *redis.PubSub) error {
	return pubsub.Unsubscribe(context.Background(), channel)
}

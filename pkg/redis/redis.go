package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisService interface {
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

// NewRedisService creates a new instance of RedisService
func New(addr string, password string, db int, ssl bool) *redisService {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		// TLSConfig: &tls.Config{
		// 	InsecureSkipVerify: false,
		// },
	})

	return &redisService{client: client}
}

// GET value from key
func (r *redisService) Get(key string) *redis.StringCmd {
	return r.client.Get(context.Background(), key)
}

// SET key value with expiration
func (r *redisService) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(context.Background(), key, value, expiration)
}

// DEL keys from Redis
func (r *redisService) Del(keys ...string) *redis.IntCmd {
	return r.client.Del(context.Background(), keys...)
}

// Set key expiration
func (r *redisService) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	return r.client.Expire(context.Background(), key, expiration)
}

// Sets
func (r *redisService) SAdd(key string, members ...interface{}) *redis.IntCmd {
	return r.client.SAdd(context.Background(), key, members)
}

// get the members of sets
func (r *redisService) SMembers(key string) *redis.StringSliceCmd {
	return r.client.SMembers(context.Background(), key)
}

// remove member from sets
func (r *redisService) SRem(key string, members ...interface{}) *redis.IntCmd {
	return r.client.SRem(context.Background(), key, members)
}

// Lists
func (r *redisService) LPush(key string, values ...interface{}) *redis.IntCmd {
	return r.client.LPush(context.Background(), key, values...)
}

func (r *redisService) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	return r.client.LRange(context.Background(), key, start, stop)
}

func (r *redisService) LPop(key string) *redis.StringCmd {
	return r.client.LPop(context.Background(), key)
}

// Hashes
func (r *redisService) HSet(key, field, value string) {
	r.client.HSet(context.Background(), key, field, value)
}

func (r *redisService) HGetAll(key string) *redis.StringStringMapCmd {
	return r.client.HGetAll(context.Background(), key)
}

// Streams
func (r *redisService) XAdd(stream string, values map[string]interface{}) {
	r.client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: stream,
		Values: values,
	})
}

func (r *redisService) XRange(stream, start, stop string) *redis.XMessageSliceCmd {
	return r.client.XRange(context.Background(), stream, start, stop)
}

// JSON
func (r *redisService) SetJson(key string, value interface{}) {
	json, _ := json.Marshal(value)
	r.client.Set(context.Background(), key, json, 0)
}

// Publish a message to a channel
func (r *redisService) Publish(channel string, message string) *redis.IntCmd {
	return r.client.Publish(context.Background(), channel, message)
}

// Subscribe to a channel and receive messages
func (r *redisService) Subscribe(channel string) *redis.PubSub {
	return r.client.Subscribe(context.Background(), channel)
}

// Remove a channel subscription
func (r *redisService) Unsubscribe(channel string, pubsub *redis.PubSub) error {
	return pubsub.Unsubscribe(context.Background(), channel)
}

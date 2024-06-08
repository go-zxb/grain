package redisx

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type IRedis interface {
	Subscribe(channel string) *redis.PubSub
	Publish(data []byte, channel string) error
	Set(key string, value interface{}, ex time.Duration)
	Get(key string) string
	Del(key string) int64
	GetInt64(key string) (int64, error)
	GetInt(key string) (int, error)
	SetInt(key string, value int64, expiration time.Duration) error
	IncrInt(key string, value int64) (int64, error)
	DecrInt(key string, value int64) (int64, error)
	GetFloat(key string) (float64, error)
	IncrFloat(key string, value float64) (float64, error)
	SetFloat(key string, value float64, expiration time.Duration) error
	GetObject(key string, v interface{}) error
	SetObject(key string, value interface{}, expiration time.Duration) error
	Incr(key string, value interface{}) (interface{}, error)
	SetEx(key string, t time.Duration)
	Exists(key string) (bool, error)
	UserSign(userID string) error
	ZRange(key string) []string
	ZAdd(key string, data interface{}) error
	SetNX(key string, value interface{}, expiration time.Duration) error
	Scan(key string, count int64) []string
	GetTTL(key string) float64
	Enqueue(key string, item interface{}) error
	Dequeue(key string, item interface{}) error
	Peek(key string, item interface{}) error
	Length(key string) (int64, error)
	Clear(key string) error
	EnqueueWithTTL(key string, item interface{}, ttl time.Duration) error
}

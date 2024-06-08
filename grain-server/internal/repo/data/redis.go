// Copyright © 2023 Grain. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-grain/grain/config"
	jsonx "github.com/go-grain/grain/pkg/encoding/json"
	redisx "github.com/go-grain/grain/pkg/redis"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

var rdb *Redis

type Redis struct {
	Client *redis.Client
	ctx    context.Context
}

func InitRedis() (client redisx.IRedis, err error) {
	conf := config.GetConfig().DataBase.Redis
	rdbClient := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		Username:     conf.UserName,
		DB:           conf.DB,
		WriteTimeout: conf.WriteTimeout,
		ReadTimeout:  conf.ReadTimeout,
	})
	_, err = rdbClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.New("Redis 连接失败: " + err.Error())
	}
	xlog.Debug("初始化Redis成功")
	rdb = &Redis{
		Client: rdbClient,
		ctx:    context.Background(),
	}
	return rdb, nil
}

func GetRedis() *Redis {
	return rdb
}

func (rs Redis) Subscribe(channel string) *redis.PubSub {
	pubs := rs.Client.Subscribe(rs.ctx, channel)
	return pubs
}

func (rs Redis) Publish(data []byte, channel string) error {
	err := rs.Client.Publish(rs.ctx, channel, data).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rs Redis) Set(key string, value interface{}, ex time.Duration) {
	err := rs.Client.Set(rs.ctx, key, value, ex*time.Second).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func (rs Redis) Get(key string) string {
	val, err := rs.Client.Get(rs.ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}

func (rs Redis) Del(key string) int64 {
	val, err := rs.Client.Del(rs.ctx, key).Result()
	if err != nil {
		return 0
	}
	return val
}

func (rs Redis) GetInt64(key string) (int64, error) {
	result, err := rs.Client.Get(rs.ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return strconv.ParseInt(result, 10, 64)
}

func (rs Redis) GetInt(key string) (int, error) {
	result, err := rs.Client.Get(rs.ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return strconv.Atoi(result)
}

func (rs Redis) SetInt(key string, value int64, expiration time.Duration) error {
	return rs.Client.Set(rs.ctx, key, strconv.FormatInt(value, 10), expiration*time.Second).Err()
}

func (rs Redis) IncrInt(key string, value int64) (int64, error) {
	return rs.Client.IncrBy(rs.ctx, key, value).Result()
}

func (rs Redis) DecrInt(key string, value int64) (int64, error) {
	return rs.Client.DecrBy(rs.ctx, key, value).Result()
}

func (rs Redis) GetFloat(key string) (float64, error) {
	result, err := rs.Client.Get(rs.ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return strconv.ParseFloat(result, 64)
}

func (rs Redis) SetFloat(key string, value float64, expiration time.Duration) error {
	return rs.Client.Set(rs.ctx, key, strconv.FormatFloat(value, 'f', -1, 64), expiration*time.Second).Err()
}

func (rs Redis) IncrFloat(key string, value float64) (float64, error) {
	return rs.Client.IncrByFloat(rs.ctx, key, value).Result()
}

func (rs Redis) GetObject(key string, v interface{}) error {
	result, err := rs.Client.Get(rs.ctx, key).Result()
	if err != nil {
		//if errors.Is(err, redis.Nil) {
		//	return nil
		//}
		return err
	}
	return jsonx.Unmarshal([]byte(result), v)
}

func (rs Redis) SetObject(key string, value any, expiration time.Duration) error {
	data := jsonx.Marshal(value)
	if data == nil {
		return errors.New("SetObject Fail")
	}
	return rs.Client.Set(rs.ctx, key, data, expiration*time.Second).Err()
}

func (rs Redis) Incr(key string, value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case int64:
		return rs.IncrInt(key, v)
	case float64:
		return rs.IncrFloat(key, v)
	default:
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}

func (rs Redis) SetEx(key string, t time.Duration) {
	rs.Client.Expire(rs.ctx, key, t*time.Second)
}

func (rs Redis) GetTTL(key string) float64 {
	ttlResult, err := rs.Client.TTL(rs.ctx, key).Result()
	if err != nil {
		// 处理错误
	}
	return ttlResult.Seconds()
}

func (rs Redis) Scan(key string, count int64) []string {
	var cursor uint64
	var keys []string
	var allKeys []string
	for {

		var err error
		keys, cursor, err = rs.Client.Scan(rs.ctx, cursor, fmt.Sprintf("%s:*", key), 10).Result()
		if err != nil {
			return nil
		}

		for _, key := range keys {
			allKeys = append(allKeys, key)
		}

		// 没有更多key了
		if cursor == 0 {
			break
		}
	}
	return allKeys
}

func (rs Redis) SetNX(key string, value interface{}, expiration time.Duration) error {
	result := rdb.Client.SetNX(rs.ctx, key, value, expiration*time.Second)
	if result.Val() {
		return nil
	} else {
		return errors.New("key已存在")
	}
}

func (rs Redis) ZAdd(key string, data interface{}) error {
	marshal := jsonx.Marshal(data)
	if marshal == nil {
		return errors.New("ZAdd操作失败")
	}
	score := float64(time.Now().Unix())
	member := marshal
	result := rdb.Client.ZAdd(rs.ctx, key, redis.Z{Score: score, Member: member})
	_, err := result.Result()
	if err != nil {
		return err
	}
	return nil
}

func (rs Redis) ZRange(key string) []string {
	var v []string
	result := rdb.Client.ZRange(rs.ctx, key, 0, -1)
	dates := result.Val()
	for _, val := range dates {
		v = append(v, val)
	}
	return v
}

func (rs Redis) UserSign(userID string) error {
	date := time.Now().Format("2006-01-02")
	score := float64(time.Now().Unix())
	member := date

	key := fmt.Sprintf("user:%s:sign", userID)
	// 判断用户是否已经签到过了
	result := rdb.Client.ZRank(rs.ctx, key, member)
	_, err := result.Result()
	if err != nil {
		// 用户未签到过，执行签到操作
		rdb.Client.ZAdd(rs.ctx, key, redis.Z{Score: score, Member: member})
		xlog.Debug("签到成功！")

		//// 查询用户签到情况
		//result := client.ZRange(rs.ctx, key, 0, -1)
		//dates := result.Val()
		//for _, date := range dates {
		//	xlog.Debug(date)
		//}
		return nil
	} else {
		// 用户已经签到过，不能重复签到
		return errors.New("您今天已经签到过了")
	}
}

func (rs Redis) Exists(key string) (bool, error) {
	exists, err := rdb.Client.Exists(rs.ctx, key).Result()
	if err != nil {
		return false, err
	}

	if exists == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (rs Redis) Enqueue(key string, item interface{}) error {
	data := jsonx.Marshal(item)
	return rs.Client.LPush(rs.ctx, key, data).Err()
}

func (rs Redis) Dequeue(key string, item interface{}) error {
	result, err := rs.Client.BRPop(rs.ctx, 0, key).Result()
	if err != nil {
		return err
	}
	if len(result) != 2 {
		return errors.New("invalid result length")
	}
	data := result[1]
	return jsonx.Unmarshal([]byte(data), item)
}

// Peek returns the first item from the queue without removing it.
func (rs Redis) Peek(key string, item interface{}) error {
	result, err := rs.Client.LIndex(rs.ctx, key, 0).Result()
	if err == redis.Nil {
		return errors.New("queue is empty")
	}

	return jsonx.Unmarshal([]byte(result), item)
}

// Length returns the number of items in the queue.
func (rs Redis) Length(key string) (int64, error) {
	return rs.Client.LLen(rs.ctx, key).Result()
}

// Clear removes all items from the queue.
func (rs Redis) Clear(key string) error {
	return rs.Client.Del(rs.ctx, key).Err()
}

// EnqueueWithTTL adds an item to the end of the queue with a TTL (time-to-live) value.
func (rs Redis) EnqueueWithTTL(key string, item interface{}, ttl time.Duration) error {
	data := jsonx.Marshal(item)
	_, err := rs.Client.Pipelined(rs.ctx, func(pipe redis.Pipeliner) error {
		pipe.LPush(rs.ctx, key, data)
		pipe.Expire(rs.ctx, key, ttl*time.Second)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

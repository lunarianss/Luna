// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package redis

import (
	"fmt"
	"sync"

	db "github.com/lunarianss/Luna/infrastructure/redis"
	"github.com/lunarianss/Luna/internal/infrastructure/options"
	"github.com/redis/go-redis/v9"
)

var (
	once     sync.Once
	RedisIns *redis.Client
)

func GetRedisIns(opt *options.RedisOptions) (*redis.Client, error) {
	if opt == nil && RedisIns == nil {
		return nil, fmt.Errorf("failed to get redis store factory")
	}

	var err error
	var redisClient redis.UniversalClient

	once.Do(func() {
		redisClient, err = db.NewRedisClusterPool(true, opt)
		RedisIns, _ = redisClient.(*redis.Client)
	})

	if RedisIns == nil || err != nil {
		return nil, fmt.Errorf("failed to get redis store factory option is %+v: error: %s", opt, err.Error())
	}
	return RedisIns, nil
}

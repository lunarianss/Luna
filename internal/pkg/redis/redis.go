package redis

import (
	"fmt"
	"sync"

	"github.com/lunarianss/Luna/internal/pkg/options"
	db "github.com/lunarianss/Luna/pkg/redis"
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

	once.Do(func() {
		RedisIns = db.NewRedisClusterPool(true, opt).(*redis.Client)
	})

	if RedisIns == nil {
		return nil, fmt.Errorf("failed to get redis store factory option is %+v", opt)
	}
	return RedisIns, nil
}

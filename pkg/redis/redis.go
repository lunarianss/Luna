// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package redis

// Package redis provides an redis implementation of the AnalyticsStorage storage interface.

import (
	"context"
	"crypto/tls"
	"strconv"
	"time"

	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/options"
	"github.com/lunarianss/Luna/pkg/errors"
	"github.com/lunarianss/Luna/pkg/log"
	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
)

// ------------------- REDIS CLUSTER STORAGE MANAGER -------------------------------

// RedisKeyPrefix defines prefix for iam analytics key.
const (
	RedisKeyPrefix      = "luna-"
	defaultRedisAddress = "127.0.0.1:6379"
)

var redisClusterSingleton redis.UniversalClient

// RedisClusterStorageManager is a storage manager that uses the redis database.
type RedisClusterStorageManager struct {
	db        redis.UniversalClient
	KeyPrefix string
	HashKeys  bool
	Config    *options.RedisOptions
}

// NewRedisClusterPool returns a redis cluster client.
func NewRedisClusterPool(forceReconnect bool, config *options.RedisOptions) (redis.UniversalClient, error) {
	if !forceReconnect {
		if redisClusterSingleton != nil {
			log.Debug("Redis pool already INITIALIZED")

			return redisClusterSingleton, nil
		}
	} else {
		if redisClusterSingleton != nil {
			redisClusterSingleton.Close()
		}
	}

	log.Debug("Creating new Redis connection pool")

	maxActive := 500
	if config.MaxActive > 0 {
		maxActive = config.MaxActive
	}

	timeout := 5 * time.Second

	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}

	var tlsConfig *tls.Config
	if config.UseSSL {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: config.SSLInsecureSkipVerify,
		}
	}

	var client redis.UniversalClient
	opts := &RedisOpts{
		MasterName:      config.MasterName,
		Addrs:           getRedisAddrs(config),
		DB:              config.Database,
		Password:        config.Password,
		PoolSize:        maxActive,
		ConnMaxIdleTime: 240 * time.Second,
		ReadTimeout:     timeout,
		WriteTimeout:    timeout,
		DialTimeout:     timeout,
		TLSConfig:       tlsConfig,
	}

	if opts.MasterName != "" {
		log.Info("--> [REDIS] Creating sentinel-backed failover client")
		client = redis.NewFailoverClient(opts.failover())
	} else if config.EnableCluster {
		log.Info("--> [REDIS] Creating cluster client")
		client = redis.NewClusterClient(opts.cluster())
	} else {
		log.Info("--> [REDIS] Creating single-node client")
		client = redis.NewClient(opts.simple())
	}

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		return nil, err
	}

	redisClusterSingleton = client
	return client, nil
}

func getRedisAddrs(config *options.RedisOptions) (addrs []string) {
	if len(config.Addrs) != 0 {
		addrs = config.Addrs
	}

	if len(addrs) == 0 && config.Port != 0 {
		addr := config.Host + ":" + strconv.Itoa(config.Port)
		addrs = append(addrs, addr)
	}

	return addrs
}

// RedisOpts is the overridden type of redis.UniversalOptions. simple() and cluster() functions are not public
// in redis library. Therefore, they are redefined in here to use in creation of new redis cluster logic.
// We don't want to use redis.NewUniversalClient() logic.
type RedisOpts redis.UniversalOptions

func (o *RedisOpts) cluster() *redis.ClusterOptions {
	if len(o.Addrs) == 0 {
		o.Addrs = []string{defaultRedisAddress}
	}

	return &redis.ClusterOptions{
		Addrs:     o.Addrs,
		OnConnect: o.OnConnect,

		Password: o.Password,

		MaxRedirects:   o.MaxRedirects,
		ReadOnly:       o.ReadOnly,
		RouteByLatency: o.RouteByLatency,
		RouteRandomly:  o.RouteRandomly,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:     o.DialTimeout,
		ReadTimeout:     o.ReadTimeout,
		WriteTimeout:    o.WriteTimeout,
		PoolSize:        o.PoolSize,
		MinIdleConns:    o.MinIdleConns,
		ConnMaxLifetime: o.ConnMaxLifetime,
		PoolTimeout:     o.PoolTimeout,
		ConnMaxIdleTime: o.ConnMaxIdleTime,

		TLSConfig: o.TLSConfig,
	}
}

func (o *RedisOpts) simple() *redis.Options {
	addr := defaultRedisAddress
	if len(o.Addrs) > 0 {
		addr = o.Addrs[0]
	}

	return &redis.Options{
		Addr:      addr,
		OnConnect: o.OnConnect,

		DB:       o.DB,
		Password: o.Password,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,

		PoolSize:        o.PoolSize,
		MinIdleConns:    o.MinIdleConns,
		ConnMaxLifetime: o.ConnMaxLifetime,
		PoolTimeout:     o.PoolTimeout,
		ConnMaxIdleTime: o.ConnMaxIdleTime,

		TLSConfig: o.TLSConfig,
	}
}

func (o *RedisOpts) failover() *redis.FailoverOptions {
	if len(o.Addrs) == 0 {
		o.Addrs = []string{"127.0.0.1:26379"}
	}

	return &redis.FailoverOptions{
		SentinelAddrs: o.Addrs,
		MasterName:    o.MasterName,
		OnConnect:     o.OnConnect,

		DB:       o.DB,
		Password: o.Password,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,

		PoolSize:        o.PoolSize,
		MinIdleConns:    o.MinIdleConns,
		ConnMaxLifetime: o.ConnMaxLifetime,
		PoolTimeout:     o.PoolTimeout,
		ConnMaxIdleTime: o.ConnMaxIdleTime,
		TLSConfig:       o.TLSConfig,
	}
}

// GetName returns the redis cluster storage manager name.
func (r *RedisClusterStorageManager) GetName() string {
	return "redis"
}

// Init initialize the redis cluster storage manager.
func (r *RedisClusterStorageManager) Init(config interface{}) error {
	r.Config = &options.RedisOptions{}
	err := mapstructure.Decode(config, &r.Config)
	if err != nil {
		log.Fatalf("Failed to decode configuration: %s", err.Error())
	}

	r.KeyPrefix = RedisKeyPrefix

	return nil
}

// Connect will establish a connection to the r.db.
func (r *RedisClusterStorageManager) Connect() bool {
	if r.db == nil {
		log.Debug("Connecting to redis cluster")
		r.db, _ = NewRedisClusterPool(false, r.Config)
		return true
	}

	log.Debug("Storage Engine already initialized...")

	// Reset it just in case
	r.db = redisClusterSingleton

	return true
}

func (r *RedisClusterStorageManager) hashKey(in string) string {
	return in
}

func (r *RedisClusterStorageManager) fixKey(keyName string) string {
	setKeyName := r.KeyPrefix + r.hashKey(keyName)

	log.Debugf("Input key was: %s", setKeyName)

	return setKeyName
}

// SetKey will create (or update) a key value in the store.
func (r *RedisClusterStorageManager) SetKey(ctx context.Context, keyName, value string, timeout int64) error {
	log.Debugf("[STORE] SET Raw key is: %s", keyName)
	log.Debugf("[STORE] Setting key: %s", r.fixKey(keyName))

	r.ensureConnection()
	err := r.db.Set(ctx, r.fixKey(keyName), value, 0).Err()

	if err != nil {
		return errors.WithCode(code.ErrRedisSetKey, err.Error())
	}

	if timeout > 0 {
		if expErr := r.SetExp(ctx, keyName, timeout); expErr != nil {
			return errors.WithCode(code.ErrRedisSetExpire, expErr.Error())
		}
	}
	return nil
}

// SetExp is used to set the expiry of a key.
func (r *RedisClusterStorageManager) SetExp(ctx context.Context, keyName string, timeout int64) error {
	if err := r.db.Expire(ctx, r.fixKey(keyName), time.Duration(timeout)*time.Second).Err(); err != nil {
		return errors.WithCode(code.ErrRedisSetExpire, err.Error())
	}
	return nil
}

func (r *RedisClusterStorageManager) ensureConnection() {
	if r.db != nil {
		// already connected
		return
	}
	log.Info("Connection dropped, reconnecting...")
	for {
		r.Connect()
		if r.db != nil {
			// reconnection worked
			return
		}
		log.Info("Reconnecting again...")
	}
}

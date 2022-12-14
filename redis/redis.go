package redis

import (
	"fmt"
	"sync"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

var (
	// Global handle to the redis
	redisp = &redisHolder{client: nil}
)

type redisHolder struct {
	client *redis.Client
	mu     sync.RWMutex
}

func init() {
	common.SetDefProp(common.PROP_REDIS_ENABLED, false)
	common.SetDefProp(common.PROP_REDIS_ADDRESS, "localhost")
	common.SetDefProp(common.PROP_REDIS_PORT, 6379)
	common.SetDefProp(common.PROP_REDIS_USERNAME, "")
	common.SetDefProp(common.PROP_REDIS_PASSWORD, "")
	common.SetDefProp(common.PROP_REDIS_DATABASE, 0)
}

/*
	Check if redis is enabled

	This func looks for following prop:

		PROP_REDIS_ENABLED
*/
func IsRedisEnabled() bool {
	return common.GetPropBool(common.PROP_REDIS_ENABLED)
}

/*
	Get Redis client

	Must call InitRedis(...) method before this method.
*/
func GetRedis() *redis.Client {
	redisp.mu.RLock()
	defer redisp.mu.RUnlock()

	if redisp.client == nil {
		panic("Redis Connection hasn't been initialized yet")
	}
	return redisp.client
}

/*
	Initialize redis client from configuration

	If redis client has been initialized, current func call will be ignored.

	This func looks for following prop:

		PROP_REDIS_ADDRESS
		PROP_REDIS_PORT
		PROP_REDIS_USERNAME
		PROP_REDIS_PASSWORD
		PROP_REDIS_DATABASE
*/
func InitRedisFromProp() *redis.Client {
	return InitRedis(
		common.GetPropStr(common.PROP_REDIS_ADDRESS),
		common.GetPropStr(common.PROP_REDIS_PORT),
		common.GetPropStr(common.PROP_REDIS_USERNAME),
		common.GetPropStr(common.PROP_REDIS_PASSWORD),
		common.GetPropInt(common.PROP_REDIS_DATABASE))
}

/*
	Initialize redis client

	If redis client has been initialized, current func call will be ignored
*/
func InitRedis(address string, port string, username string, password string, db int) *redis.Client {
	if IsRedisClientInitialized() {
		return GetRedis()
	}

	redisp.mu.Lock()
	defer redisp.mu.Unlock()

	if redisp.client != nil {
		return redisp.client
	}

	logrus.Infof("Connecting to redis '%v:%v', database: %v", address, port, db)
	var rdb *redis.Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", address, port),
		Password: password,
		DB:       db,
	})

	logrus.Info("Redis Handle initialized")
	redisp.client = rdb
	return rdb
}

// Check whether redis client is initialized
func IsRedisClientInitialized() bool {
	redisp.mu.RLock()
	defer redisp.mu.RUnlock()
	return redisp.client != nil
}

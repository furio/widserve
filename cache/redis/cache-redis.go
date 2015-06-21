package redis

import (
	"time"
	"strconv"
	"sync"
	"fmt"
	"gopkg.in/redis.v3"
)

type RedisCache struct {
	instance *redis.Client
}

var init_ctx sync.Once


func (this RedisCache) Init(config map[string]string) {
	init_ctx.Do( func () {
		database, _ := strconv.ParseInt(config["database"], 10, 64)
		poolSize, _ := strconv.ParseInt(config["poolSize"], 10, 64)

		this.instance = redis.NewClient(&redis.Options{
			Addr:     config["address"],
			Password: config["password"],
			DB:       database,
			PoolSize: int(poolSize),
		})
	})
}

func (this RedisCache) Get(key string) (interface{},bool) {
	status := this.instance.Get(key);

	if ( status.Err() != nil || len(status.Val()) == 0) {
		return nil, false;
	}

	return status.Val(), true
}

func (this RedisCache) Set(key string, value interface{}, timeout time.Duration) bool {
	if str, ok := value.(string); ok {
		status := this.instance.Set(key,str,timeout)
		return status.Err() == nil
	} else if str, ok := value.(fmt.Stringer); ok {
		status := this.instance.Set(key,str.String(),timeout)
		return status.Err() == nil
	} else {
		return false
	}
}

func (this RedisCache) Delete(key string) bool {
	status := this.instance.Del(key)
	return status.Err() == nil && status.Val() == 1
}
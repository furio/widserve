package local

import (
	"github.com/pmylund/go-cache"
	"time"
	"strconv"
	"sync"
)

type LocalCache struct {
	instance *cache.Cache
}

var init_ctx sync.Once


func (this LocalCache) Init(config map[string]string) {
	init_ctx.Do( func () {
			defaultExpiration, _ := strconv.ParseInt(config["defaultExpiration"], 10, 64)
			purgeTime, _ := strconv.ParseInt(config["purgeTime"], 10, 64)

			this.instance = cache.New(time.Duration(defaultExpiration)*time.Second, time.Duration(purgeTime)*time.Second)
	})
}

func (this LocalCache) Get(key string) (interface{},bool) {
	return this.instance.Get(key);
}

func (this LocalCache) Set(key string, value interface{}, timeout time.Duration)  {
	this.instance.Set(key,value,timeout)
}

func (this LocalCache) Delete(key string) {
	this.instance.Delete(key)
}
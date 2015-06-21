package db

import (
	"time"
)

type Widget struct {
	WidgetID	string	`db:"widget_id,size:255"`
	ApiKey		string	`db:"api_key,size:255"`
	ApiPath		string	`db:"api_key,size:1024"`
	Created		int64	`db:"created_at"`
	CacheElapse	int64	`db:"cache_elapse"`
	LastCache	int64	`db:"last_cache_check"`
}

func NewWidget(uid string, apiKey string, apiPath string, cacheElapse int64) Widget {
	now := time.Now().UnixNano()

	return Widget{
		WidgetID: uid,
		Created: now,
		ApiKey: apiKey,
		ApiPath: apiPath,
		CacheElapse: cacheElapse,
		LastCache: now - cacheElapse,
	}
}
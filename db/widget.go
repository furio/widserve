package db

import (
	"time"
)

type Widget struct {
	WidgetID	string	// `db:"widget_id,size:255"`
	ApiKey		string	// `db:"api_key,size:255"`
	ApiPath		string	// `db:"api_path,size:1024"`
	Created		uint64	// `db:"created_at"`
	CacheElapse	uint32	// `db:"cache_elapse"`
	NextCheck	uint64	// `db:"next_cache_check"`
}

func NewWidget(uid string, apiKey string, apiPath string, cacheElapse uint32) Widget {
	now := uint64( time.Now().UnixNano() )

	return Widget{
		WidgetID: uid,
		Created: now,
		ApiKey: apiKey,
		ApiPath: apiPath,
		CacheElapse: cacheElapse,
		NextCheck: now + uint64(cacheElapse),
	}
}
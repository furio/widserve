package db

import (
	"fmt"
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

func (this Widget) String() string{
	return fmt.Sprintf("{Widget Id:%s, ApiKey:%s, ApiPath:%s, Created:%d, CacheElapse:%d, NextCheck:%d}",
		this.WidgetID, this.ApiKey, this.ApiPath, this.Created, this.CacheElapse, this.NextCheck)
}

func newWidget(uid string, apiKey string, apiPath string, cacheElapse uint32) Widget {
	now := uint64( time.Now().Unix() )

	return Widget{
		WidgetID: uid,
		Created: now,
		ApiKey: apiKey,
		ApiPath: apiPath,
		CacheElapse: cacheElapse,
		NextCheck: now + uint64(cacheElapse),
	}
}
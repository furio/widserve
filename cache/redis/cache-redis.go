package redis

import (
	"log"

	"gopkg.in/redis.v3"
)

var _ = log.Logger{}
var _ = redis.Nil

/*
var redisCli = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})
*/
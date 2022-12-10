package config

import (
	"github.com/chiyoi/neko03/pkg/logs"
	"github.com/chiyoi/trinity/internal/configs"
	"github.com/go-redis/redis/v8"
)

var Config = map[string]any{
	"ServiceURL": "http://maru/",

	"TrinityURL":     "http://trinity/",
	"OnebotURL":      "http://gocq/",
	"OnebotEventURL": "ws://gocq:8080/",

	"RedisOptions": &redis.Options{
		Addr:     "redis-18080.c56.east-us.azure.cloud.redislabs.com:18080",
		Username: "maru",
		Password: "Neko03Maru@redis",
	},
	"RedisKeyUsersLoggedIn": "maru:usersLoggedIn",
	"RedisKeyNekoMap":       "maru:nekoMap",
	"RedisKeyListeners":     "trinity:listeners",
}

func GetErr[T any](key string) (T, error) {
	return configs.Get[T](Config, key)
}

func Get[T any](key string) (val T) {
	val, err := GetErr[T](key)
	if err != nil {
		logs.Panic(err)
	}
	return
}

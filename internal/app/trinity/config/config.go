package config

import (
	"github.com/chiyoi/neko03/pkg/logs"
	"github.com/chiyoi/trinity"
	"github.com/chiyoi/trinity/internal/configs"
)

func GetErr[T any](key string) (T, error) {
	return configs.Get[T](trinity.TrinityConfig, key)
}

func Get[T any](key string) (val T) {
	val, err := GetErr[T](key)
	if err != nil {
		logs.Panic(err)
	}
	return
}

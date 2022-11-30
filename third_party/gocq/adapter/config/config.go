package config

import (
	"github.com/chiyoi/trinity"
	"github.com/chiyoi/trinity/internal/configs"
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

func GetErr[T any](key string) (T, error) {
	return configs.Get[T](trinity.GocqConfig, key)
}

func Get[T any](key string) (val T) {
	val, err := GetErr[T](key)
	if err != nil {
		logs.Fatal(err)
	}
	return
}

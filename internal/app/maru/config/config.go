package config

import (
	"github.com/chiyoi/trinity"
	"github.com/chiyoi/trinity/internal/configs"
)

var cfg = trinity.MaruConfig

func Get[T any](key string) (a T, err error) {
	return configs.Get[T](cfg, key)
}

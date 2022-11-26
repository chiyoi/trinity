package maru

import (
	"github.com/chiyoi/trinity"
	"github.com/chiyoi/trinity/internal/configs"
)

var cfg = trinity.MaruConfig

func GetConfig[T any](key string) (a T, err error) {
	return configs.Get[T](cfg, key)
}

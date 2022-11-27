package maru

import (
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

var (
	serviceUrl string

	trinityUrl string
)

func init() {
	var err error
	if serviceUrl, err = GetConfig[string]("ServiceURL"); err != nil {
		logs.Fatal("maru:", err)
	}

	if trinityUrl, err = GetConfig[string]("TrinityURL"); err != nil {
		logs.Fatal("maru:", err)
	}
}

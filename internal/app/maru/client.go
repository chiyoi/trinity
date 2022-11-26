package maru

import (
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

var trinityUrl string

func init() {
	var err error
	if trinityUrl, err = GetConfig[string]("TrinityURL"); err != nil {
		logs.Fatal("maru:", trinityUrl)
	}
}

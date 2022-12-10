package aira

import (
	"time"

	"github.com/chiyoi/neko03/pkg/logs"
	"github.com/chiyoi/trinity/internal/app/aira/config"
	"github.com/chiyoi/trinity/internal/app/aira/db"
)

func Heartbeat() {
	serviceURL := config.Get[string]("ServiceURL")
	if err := db.RegisterListener(serviceURL); err != nil {
		logs.Panic(err)
	}
	defer func() {
		if err := db.RemoveListener(serviceURL); err != nil {
			logs.Error(err)
		}
	}()
	for {
		time.Sleep(time.Second * 10)
		ok, err := db.CheckListener(serviceURL)
		if err != nil {
			logs.Panic(err)
		}
		if !ok {
			if err := db.RegisterListener(serviceURL); err != nil {
				logs.Panic(err)
			}
		}
	}
}

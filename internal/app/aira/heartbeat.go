package aira

import (
	"time"

	"github.com/chiyoi/trinity/internal/app/aira/client"
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

func Heartbeat() {
	logPrefix := "heartbeat:"
	for {
		time.Sleep(time.Second * 10)
		if ok, err := client.CheckListener(); err != nil {
			logs.Warning(logPrefix, err)
			if !ok {
				client.RegisterListener()
			}
		}
	}
}

package main

import (
	"github.com/chiyoi/trinity/internal/app/aira"
	"github.com/chiyoi/trinity/internal/app/aira/client"
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

func main() {
	rdb, err := aira.OpenRedis()
	if err != nil {
		logs.Fatal("aira:", err)
	}
	if err = client.RegisterListener(rdb); err != nil {
		logs.Fatal("aira:", err)
	}
	timestampChannel := make(chan int64, 1)
	go func() {
		if err := client.EventSynchronizer(timestampChannel, rdb); err != nil {
			logs.Fatal("aira:", err)
		}
	}()
	srv := aira.Server(timestampChannel)
	go aira.StartSrv(srv)
	defer aira.StopSrv(srv)
	select {}
}

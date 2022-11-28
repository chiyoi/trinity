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
	client.RegisterListener(rdb)
	chanTimestamp := make(chan int64, 1)
	go func() {
		if err := client.EventSynchronizer(chanTimestamp, rdb); err != nil {
			logs.Fatal("aira:", err)
		}
	}()
	srv := aira.Server(chanTimestamp)
	go aira.StartSrv(srv)
	defer aira.StopSrv(srv)
	select {}
}

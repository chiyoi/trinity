package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/chiyoi/trinity/internal/app/aira"
	"github.com/chiyoi/trinity/internal/app/aira/client"
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

func main() {
	rdb, err := aira.OpenRedis()
	if err != nil {
		logs.Fatal(err)
	}
	if err = client.RegisterListener(rdb); err != nil {
		logs.Fatal(err)
	}
	defer func() {
		if err := client.RemoveListener(rdb); err != nil {
			logs.Error(err)
		}
	}()
	timestampChannel := make(chan int64, 1)
	go func() {
		if err := client.EventSynchronizer(timestampChannel, rdb); err != nil {
			logs.Fatal("aira:", err)
		}
	}()
	srv := aira.Server(timestampChannel)
	go aira.StartSrv(srv)
	defer aira.StopSrv(srv)

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGTERM)
	<-term
	logs.Info("terminate")
}

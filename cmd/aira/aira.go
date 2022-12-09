package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/chiyoi/trinity/internal/app/aira"
	"github.com/chiyoi/trinity/internal/app/aira/client"
	"github.com/chiyoi/trinity/internal/app/aira/db"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
)

func main() {
	rdb, err := db.OpenRedis()
	if err != nil {
		logs.Fatal(err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			logs.Error(err)
		}
	}()
	db.SetDB(rdb)

	if err = client.RegisterListener(); err != nil {
		logs.Fatal(err)
	}
	defer func() {
		if err := client.RemoveListener(); err != nil {
			logs.Error(err)
		}
	}()
	go aira.Heartbeat()

	srv := aira.Server()
	go atmt.StartSrv(srv)
	defer atmt.StopSrv(srv)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop
	logs.Info("stop:", sig)
}

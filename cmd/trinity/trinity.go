package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/chiyoi/trinity/internal/app/trinity"
	"github.com/chiyoi/trinity/internal/app/trinity/db"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
)

func main() {
	mongodb, err := db.OpenMongo()
	if err != nil {
		logs.Fatal(err)
	}
	rdb, err := db.OpenRedis()
	if err != nil {
		logs.Fatal(err)
	}
	db.SetDB(rdb, mongodb)

	srv := trinity.Server()
	go atmt.StartSrv(srv)
	defer atmt.StopSrv(srv)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop
	logs.Info("stop:", sig)
}

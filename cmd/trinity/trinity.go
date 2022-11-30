package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/chiyoi/neko03/pkg/neko"
	"github.com/chiyoi/trinity/internal/app/trinity"
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

func main() {
	mongodb, err := trinity.OpenMongo()
	if err != nil {
		logs.Fatal(err)
	}
	rdb, err := trinity.OpenRedis()
	if err != nil {
		logs.Fatal(err)
	}

	srv := trinity.Server(mongodb, rdb)
	go neko.StartSrv(srv, false)
	defer neko.StopSrv(srv)

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGTERM)
	<-term
	logs.Info("terminate")
}

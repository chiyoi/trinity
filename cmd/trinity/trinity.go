package main

import (
	"github.com/chiyoi/trinity/internal/app/trinity"
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

func main() {
	mongodb, err := trinity.OpenMongo()
	if err != nil {
		logs.Fatal("trinity:", err)
	}
	rdb, err := trinity.OpenRedis()
	if err != nil {
		logs.Fatal("trinity", err)
	}

	srv := trinity.Server(mongodb, rdb)
	go trinity.StartSrv(srv)
	defer trinity.StopSrv(srv)
	select {}
}

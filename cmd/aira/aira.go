package main

import (
	"github.com/chiyoi/neko03/pkg/logs"
	"github.com/chiyoi/neko03/pkg/neko"
	"github.com/chiyoi/trinity/internal/app/aira"
	"github.com/chiyoi/trinity/internal/app/aira/db"
	"github.com/chiyoi/trinity/pkg/atmt"
)

func main() {
	rdb, err := db.OpenRedis()
	if err != nil {
		logs.Panic(err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			logs.Error(err)
		}
	}()
	db.SetDB(rdb)
	go aira.Heartbeat()

	srv := aira.Server()
	go atmt.StartSrv(srv)
	defer atmt.StopSrv(srv)
	neko.BlockToStop()
}

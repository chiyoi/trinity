package main

import (
	"github.com/chiyoi/neko03/pkg/logs"
	"github.com/chiyoi/neko03/pkg/neko"
	"github.com/chiyoi/trinity/internal/app/trinity"
	"github.com/chiyoi/trinity/internal/app/trinity/db"
	"github.com/chiyoi/trinity/pkg/atmt"
)

func main() {
	mongodb, err := db.OpenMongo()
	if err != nil {
		logs.Panic(err)
	}
	rdb := db.OpenRedis()
	db.SetDB(rdb, mongodb)

	srv := trinity.Server()
	go atmt.StartSrv(srv)
	defer atmt.StopSrv(srv)
	neko.BlockToStop()
}

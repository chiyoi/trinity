package main

import (
	"github.com/chiyoi/trinity/internal/app/maru"
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

func main() {
	rdb, err := maru.OpenRedis()
	if err != nil {
		logs.Fatal("maru:", err)
	}

	if err = maru.RegisterListener(rdb); err != nil {
		logs.Fatal("maru:", err)
	}
	srv := maru.Server(rdb)
	serveOnebot := maru.OnebotServer(rdb)

	ws, err := maru.DialOnebotEventServer()
	if err != nil {
		logs.Fatal("maru:", err)
	}
	go serveOnebot(ws)
	go maru.StartSrv(srv)
	defer maru.StopSrv(srv)
	select {}
}

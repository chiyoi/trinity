package main

import (
	"github.com/chiyoi/trinity/internal/app/maru"
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

func main() {
	rdb, err := maru.OpenRedis()
	if err != nil {
		logs.Fatal("trinity", err)
	}
	ch := make(chan any, 10)

	onebotSrv, atmtSrv := maru.OnebotServer(rdb, ch), maru.AtmtServer(ch)
	go maru.StartOnebotSrv(onebotSrv)
	defer maru.StopOnebotSrv(onebotSrv)
	go maru.StartAtmtSrv(atmtSrv)
	defer maru.StopAtmtSrv(atmtSrv)
	select {}
}

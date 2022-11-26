package aira

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chiyoi/trinity/internal/app/aira/app"
	"github.com/chiyoi/trinity/internal/app/aira/app/aira"
	"github.com/chiyoi/trinity/internal/app/aira/app/atri"
	"github.com/chiyoi/trinity/internal/app/aira/app/eroira"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	atmt2 "github.com/chiyoi/trinity/pkg/atmt"
)

func Server(chanTimestamp chan<- int64) *atmt2.Server {
	mux := atmt2.NewServeMux()
	mux.Handle(aira.Aira())
	mux.Handle(atri.Atri())
	mux.Handle(eroira.Eroira())

	return &atmt2.Server{
		Addr:    ":http",
		Handler: app.LogEvent(mux, chanTimestamp),
	}
}

func StartSrv(srv *atmt2.Server) {
	logs.Info("aira: listening", srv.Addr)
	logs.Info("アトリ、起動！")
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		logs.Error(err)
		return
	}
	logs.Info(fmt.Sprintf("aria: server at %s closed.", srv.Addr))
}

func StopSrv(srv *atmt2.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logs.Error(err)
		return
	}
}

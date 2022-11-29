package aira

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chiyoi/trinity/internal/app/aira/handlers"
	"github.com/chiyoi/trinity/internal/app/aira/handlers/aira"
	"github.com/chiyoi/trinity/internal/app/aira/handlers/atri"
	"github.com/chiyoi/trinity/internal/app/aira/handlers/eroira"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
)

func Server(chanTimestamp chan<- int64) *atmt.Server {
	mux := atmt.NewServeMux()
	mux.Handle(aira.Aira())
	mux.Handle(atri.Atri())
	mux.Handle(eroira.Eroira())

	return &atmt.Server{
		Addr:    ":http",
		Handler: handlers.LogEvent(mux, chanTimestamp),
		ErrorCallback: map[int]func(w http.ResponseWriter, err error){
			http.StatusInternalServerError: func(w http.ResponseWriter, err error) {
				logs.Error("aira:", err)
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
			},
			http.StatusBadRequest: func(w http.ResponseWriter, err error) {
				logs.Warning("aira:", err)
				http.Error(w, "400 bad request", http.StatusBadRequest)
			},
		},
	}
}

func StartSrv(srv *atmt.Server) {
	logs.Info("アトリ、起動！")
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		logs.Error("aira:", err)
		return
	}
	logs.Info(fmt.Sprintf("aria: server at %s closed.", srv.Addr))
}

func StopSrv(srv *atmt.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logs.Error(err)
		return
	}
}

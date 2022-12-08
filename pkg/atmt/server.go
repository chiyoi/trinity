package atmt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/chiyoi/neko03/pkg/neko"
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

type Server struct {
	Addr    string
	Handler Handler
	httpSrv *http.Server
}

func (srv *Server) h() Handler {
	if srv.Handler == nil {
		return DefaultServeMux
	}
	return srv.Handler
}

func (srv *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logs.Error("trin:", err)
		neko.InternalServerError(w)
		return
	}
	var post Message
	if err = json.Unmarshal(data, &post); err != nil {
		logs.Warning("trin:", err)
		neko.BadRequest(w)
		return
	}

	var resp Message
	defer func() {
		data, err := json.Marshal(resp)
		if err != nil {
			logs.Error("trin:", err)
			neko.InternalServerError(w)
			return
		}
		if _, err = w.Write(data); err != nil {
			logs.Error("trin:", err)
			return
		}
	}()
	srv.h().ServeMessage(&resp, post)
}

func (srv *Server) ListenAndServe() (err error) {
	srv.httpSrv = &http.Server{
		Addr: srv.Addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srv.handleHTTP(w, r)
		}),
	}
	return srv.httpSrv.ListenAndServe()
}

func (srv *Server) Shutdown(ctx context.Context) (err error) {
	return srv.httpSrv.Shutdown(ctx)
}

func StartSrv(srv *Server) {
	logs.Info("listening", srv.Addr)
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		logs.Error(err)
		return
	}
	logs.Info(fmt.Sprintf("server at %s closed.", srv.Addr))
}

func StopSrv(srv *Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logs.Error(err)
		return
	}
}

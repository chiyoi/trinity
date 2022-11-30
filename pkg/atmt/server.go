package atmt

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/chiyoi/neko03/pkg/neko"
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

type Server struct {
	Addr string

	Handler Handler

	httpSrv *http.Server
}

func (srv *Server) ListenAndServe() (err error) {
	h := srv.Handler
	if h == nil {
		h = DefaultServeMux
	}

	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			logs.Error("atmt:", err)
			neko.InternalServerError(w)
			return
		}
		var req Request
		if err = json.Unmarshal(data, &req); err != nil {
			logs.Warning("atmt:", err)
			neko.BadRequest(w)
			return
		}

		ev := Event{
			time.Unix(req.Time, 0),
			req.User,
			req.MessageId,
			req.Message,
		}
		go h.ServeEvent(ev)
	}
	srv.httpSrv = &http.Server{
		Addr:    srv.Addr,
		Handler: http.HandlerFunc(httpHandler),
	}

	return srv.httpSrv.ListenAndServe()
}

func (srv *Server) Shutdown(ctx context.Context) (err error) {
	return srv.httpSrv.Shutdown(ctx)
}

func ListenAndServe(addr string, handler Handler) error {
	srv := &Server{
		Addr:    addr,
		Handler: handler,
	}
	return srv.ListenAndServe()
}

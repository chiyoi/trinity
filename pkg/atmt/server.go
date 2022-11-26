package atmt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/chiyoi/trinity/internal/pkg/logs"
)

type Server struct {
	Addr string

	Handler Handler

	httpSrv *http.Server
}

func internalServerErrorCallback(w http.ResponseWriter, err error) {
	logs.Error("atmt:", err)
	http.Error(w, "500 internal server error", http.StatusInternalServerError)
}

func badRequestCallback(w http.ResponseWriter, err error) {
	logs.Warning("atmt:", err)
	http.Error(w, "400 bad request", http.StatusBadRequest)
}

func (srv *Server) ListenAndServe() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("atmt: %w", err)
		}
	}()

	h := srv.Handler
	if h == nil {
		h = DefaultServeMux
	}

	srv.httpSrv.Addr = srv.Addr
	srv.httpSrv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			internalServerErrorCallback(w, err)
			return
		}
		var req Request
		if err = json.Unmarshal(data, &req); err != nil {
			badRequestCallback(w, err)
			return
		}

		ev := Event{
			time.Unix(int64(req.Time), 0),
			req.User,
			req.MessageId,
			req.Message,
		}
		go h.ServeEvent(ev)
	})

	return srv.httpSrv.ListenAndServe()
}

func (srv *Server) Shutdown(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("atmt: %w", err)
		}
	}()

	return srv.httpSrv.Shutdown(ctx)
}

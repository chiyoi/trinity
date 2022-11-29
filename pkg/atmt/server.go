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

	Handler       Handler
	ErrorCallback map[int]func(w http.ResponseWriter, err error)

	httpSrv *http.Server
}

var defaultErrorCallback = map[int]func(w http.ResponseWriter, err error){
	http.StatusInternalServerError: func(w http.ResponseWriter, err error) {
		logs.Error("atmt:", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
	},
	http.StatusBadRequest: func(w http.ResponseWriter, err error) {
		logs.Warning("atmt:", err)
		http.Error(w, "400 bad request", http.StatusBadRequest)
	},
}

func (srv *Server) ListenAndServe() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("atmt: %w", err)
		}
	}()
	for k, cb := range defaultErrorCallback {
		if _, ok := srv.ErrorCallback[k]; !ok {
			srv.ErrorCallback[k] = cb
		}
	}

	h := srv.Handler
	if h == nil {
		h = DefaultServeMux
	}

	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			srv.ErrorCallback[http.StatusInternalServerError](w, err)
			return
		}
		var req Request
		if err = json.Unmarshal(data, &req); err != nil {
			srv.ErrorCallback[http.StatusBadRequest](w, err)
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
	defer func() {
		if err != nil {
			err = fmt.Errorf("atmt: %w", err)
		}
	}()

	return srv.httpSrv.Shutdown(ctx)
}

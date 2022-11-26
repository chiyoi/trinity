package websocket

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type ServerConfig struct {
	// Protocol contains protocols the server support.
	// May choose one of them for each connect request.
	Protocol []string
	// Extensions contains extensions the server support.
	// for each connection request, server should support all the client extensions.
	Extensions []string

	mu sync.Mutex
}

// Hijack hijacks an existing http connection and returns a websocket connection
func (srv *ServerConfig) Hijack(w http.ResponseWriter, r *http.Request) (ws WebSocket, err error) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	defer func() {
		if err != nil {
			err = fmt.Errorf("websocket: hijack error: %w", err)
		}
	}()
	hi, ok := w.(http.Hijacker)
	if !ok {
		err = errors.New("cannot take over connection")
		return
	}

	for k, v := range map[string]string{"Upgrade": "websocket", "Connection": "Upgrade", "Sec-WebSocket-Version": "13"} {
		if got := r.Header.Get(k); got != v {
			err = fmt.Errorf("request header |%s| error: got %s, expect %s", k, got, v)
			return
		}
	}
	for _, k := range []string{"Sec-WebSocket-Key"} {
		if r.Header.Get(k) == "" {
			err = fmt.Errorf("missing request header |%s|", k)
			return
		}
	}

	h := w.Header()
	h.Set("Upgrade", "websocket")
	h.Set("Connection", "Upgrade")
	h.Set("Sec-WebSocket-Accept", challengeKey(r.Header.Get("Sec-WebSocket-Key")))

	eReq := map[string]bool{}
	for _, e := range r.Header.Values("Sec-WebSocket-Extensions") {
		eReq[e] = true
	}
	for _, e := range srv.Extensions {
		if eReq[e] {
			h.Add("Sec-WebSocket-Extensions", e)
		}
	}

	pResp := map[string]bool{}
	for _, p := range srv.Protocol {
		pResp[p] = true
	}
	for _, p := range r.Header.Values("Sec-WebSocket-Protocol") {
		if pResp[p] {
			h.Set("Sec-WebSocket-Protocol", p)
			break
		}
	}
	w.WriteHeader(http.StatusSwitchingProtocols)

	conn, _, err := hi.Hijack()
	if err != nil {
		err = fmt.Errorf("hijacker error: %w", err)
		return
	}

	ws = &webSocket{
		conn:  conn,
		srv:   true,
		state: OPEN,
		ext:   h.Values("Sec-WebSocket-Extensions"),
		pro:   h.Get("Sec-WebSocket-Protocol"),
		d:     map[Opcode][]byte{},
	}
	return
}

var defaultServerConfig ServerConfig
var DefaultServerConfig = &defaultServerConfig

func Hijack(w http.ResponseWriter, r *http.Request) (WebSocket, error) {
	return DefaultServerConfig.Hijack(w, r)
}

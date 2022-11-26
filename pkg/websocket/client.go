package websocket

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	urlpkg "net/url"
	"sync"
)

type NetDialer interface {
	Dial(network string, addr string) (net.Conn, error)
	DialContext(ctx context.Context, network string, addr string) (net.Conn, error)
}

// Dialer contains the config to start a websocket connection.
// It stores some intermediate variables during connecting,
// so it's better not to call it concurrently.
type Dialer struct {
	// Origin should be contained for a web browser.
	Origin string
	// Protocol is the sub-protocols wish to use.
	// server may choose one or none from them.
	Protocol []string
	// Extensions are websocket extensions to use.
	// server should support all of them.
	Extensions []string

	// Header contains original headers in the handshake request.
	Header http.Header
	// Body is the body for handshake request. Usually empty.
	Body io.ReadCloser

	// Dialer is the net dialer for connect,
	// tls.Dialer will be used for "wss" scheme and net.Dialer for "ws".
	// Set to a tls.Dialer with expected tls.Config if necessary.
	Dialer NetDialer

	mu  sync.Mutex
	ctx context.Context
	url *urlpkg.URL
}

// DialContext starts a websocket connection to the url using provided context,
// returns a WebSocket connection.
// The provided context must be non-nil.
func (d *Dialer) DialContext(ctx context.Context, url string) (ws WebSocket, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("websocket: dial error: %w", err)
		}
	}()
	if ctx == nil {
		err = errors.New("nil context")
		return
	}
	d.ctx = ctx

	u, err := urlpkg.Parse(url)
	if err != nil {
		return
	}
	fillHost(u)

	if err = VerifyURI(u.String()); err != nil {
		return
	}
	d.url = u

	nd := d.Dialer
	if nd == nil {
		if u.Scheme == "wss" {
			nd = new(tls.Dialer)
		} else {
			nd = new(net.Dialer)
		}
	}

	conn, err := nd.DialContext(ctx, "tcp", u.Host)
	if err != nil {
		return
	}

	if ws, err = d.handshake(conn); err != nil {
		return
	}

	return
}

func fillHost(url *urlpkg.URL) {
	if url.Scheme == "" {
		url.Scheme = "ws"
	}
	if p, ok := PortMap[url.Scheme]; ok {
		hostname, port := url.Hostname(), url.Port()
		if hostname == "" {
			hostname = "localhost"
		}
		if port == "" {
			port = p
		}
		url.Host = hostname + ":" + port
	}
}

// Dial calls DialContext with background context.
func (d *Dialer) Dial(url string) (ws WebSocket, err error) {
	return d.DialContext(context.Background(), url)
}

func (d *Dialer) handshake(conn net.Conn) (WebSocket, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	wsKey, accKey := genWsKeyPair()

	h := d.Header.Clone()
	if h == nil {
		h = http.Header{}
	}
	h.Set("Upgrade", "websocket")
	h.Set("Connection", "Upgrade")
	h.Set("Sec-WebSocket-Key", wsKey)
	h.Set("Sec-WebSocket-Version", "13")
	if d.Origin != "" {
		h.Set("Origin", d.Origin)
	}
	h.Del("Sec-WebSocket-Protocol")
	for _, p := range d.Protocol {
		h.Add("Sec-WebSocket-Protocol", p)
	}
	h.Del("Sec-WebSocket-Extensions")
	for _, e := range d.Extensions {
		h.Add("Sec-WebSocket-Extensions", e)
	}

	req := &http.Request{
		Method:           "GET",
		URL:              d.url,
		Proto:            "HTTP/1.1",
		ProtoMajor:       1,
		ProtoMinor:       1,
		Host:             d.url.Host,
		Header:           h,
		Body:             d.Body,
		TransferEncoding: []string{"identity"},
	}
	req = req.WithContext(d.ctx)

	err := req.Write(conn)
	if err != nil {
		return nil, fmt.Errorf("request write error: %w", err)
	}

	resp, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		return nil, fmt.Errorf("response read error: %w", err)
	}

	if resp.StatusCode != http.StatusSwitchingProtocols {
		return nil, errors.New("response status error: " + resp.Status)
	}

	for k, v := range map[string]string{"Upgrade": "websocket", "Connection": "Upgrade", "Sec-WebSocket-Accept": accKey} {
		if got := resp.Header.Get(k); got != v {
			return nil, fmt.Errorf("request header |%s| error: got %s, expect %s", k, got, v)
		}
	}

	eReq := make(map[string]bool)
	for _, e := range d.Extensions {
		eReq[e] = true
	}
	for _, e := range resp.Header.Values("Sec-WebSocket-Extensions") {
		if !eReq[e] {
			return nil, errors.New("response with unsupported extension: " + e)
		}
	}

	pResp := resp.Header.Get("Sec-WebSocket-Protocol")
	if !func() bool {
		if pResp == "" {
			return req.Header.Get("Sec-WebSocket-Protocol") == ""
		}
		for _, p := range req.Header.Values("Sec-WebSocket-Protocol") {
			if p == pResp {
				return true
			}
		}
		return false
	}() {
		return nil, errors.New("response protocol not supported: " + pResp)
	}

	return &webSocket{
		conn:  conn,
		state: OPEN,
		ext:   resp.Header.Values("Sec-WebSocket-Extensions"),
		pro:   resp.Header.Get("Sec-WebSocket-Protocol"),
		d:     map[Opcode][]byte{},
	}, nil
}

var defaultDialer Dialer
var DefaultDialer = &defaultDialer

// Dial connects to the url using DefaultDialer.
func Dial(url string) (ws WebSocket, err error) {
	return DefaultDialer.Dial(url)
}

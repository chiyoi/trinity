package websocket

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	ws, err := Dial("ws://localhost:12394/")
	if err != nil {
		t.Error(err)
	}
	t.Log(ws.State())

	ws, err = Dial("wss://www.neko03.com/")
	if err != nil {
		t.Error(err)
	}
	t.Log(ws.State())

	c := Dialer{
		Origin:     "http://localhost/",
		Protocol:   nil,
		Extensions: nil,
		Body:       io.NopCloser(strings.NewReader("neko nyan nyan")),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	ws, err = c.DialContext(ctx, "ws://localhost:12394/")
	if err != nil {
		t.Error(err)
	}
	t.Log(ws.State())
}

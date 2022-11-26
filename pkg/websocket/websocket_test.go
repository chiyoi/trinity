package websocket

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	Error = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix).Panicln
	Info  = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix).Println
)

func TestPing(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws, err := Hijack(w, r)
		if err != nil {
			Error(err)
			return
		}
		defer func(ws WebSocket) {
			err = ws.Close()
			if err != nil {
				Info("server close:", err)
			}
		}(ws)
		Info("server:", ws.State())

		for i := 2; i > 0; i-- {
			time.Sleep(time.Second)
			var resp []byte
			resp, err = ws.Ping()
			if err != nil {
				Error("server error:", err)
				return
			}
			Info("server received:", string(resp))
		}
	})
	go http.ListenAndServe("", nil)
	time.Sleep(time.Second * 2)

	ws, err := Dial("")
	if err != nil {
		Error(err)
	}
	defer func(ws WebSocket) {
		err = ws.Close()
		if err != nil {
			Info("client close:", err)
		}
	}(ws)
	Info("client:", ws.State())

	for {
		var rev []byte
		rev, err = ws.Recv()
		if err != nil {
			if _, ok := err.(ConnectionCloseError); ok {
				Info("client close:", err)
				return
			}
			Error("client error:", err)
			break
		}
		Info("client received:", string(rev))
	}

	time.Sleep(time.Second)
}

func TestMessage(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws, err := Hijack(w, r)
		if err != nil {
			Error(err)
		}
		Info("server:", ws.State())

		if err = ws.SendText("Hello."); err != nil {
			Error("server error:", err)
			return
		}

		resp, err := ws.Recv()
		if err != nil {
			Error("server error:", err)
			return
		}
		Info("server received:", string(resp))

		var msg = strings.Repeat("nacho ", 10000)
		if err = ws.SendText(msg); err != nil {
			Error("server error:", err)
			return
		}
	})
	go http.ListenAndServe("", nil)
	time.Sleep(time.Second * 2)

	ws, err := Dial("")
	if err != nil {
		Error(err)
	}
	Info("server:", ws.State())

	var rev []byte
	if rev, err = ws.Recv(); err != nil {
		Error("client error:", err)
		return
	}
	Info("client received:", string(rev))

	if err = ws.SendText("Hi."); err != nil {
		Error("client error:", err)
		return
	}

	if rev, err = ws.Recv(); err != nil {
		Error("client error:", err)
		return
	}
	Info("client received:", string(rev))
	time.Sleep(time.Second)
}

func TestAbnormalClosure(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = Hijack(w, r)
		select {}
	})
	go http.ListenAndServe("", nil)
	time.Sleep(time.Second * 2)

	ws, err := Dial("")
	if err != nil {
		Error(err)
	}

	go ws.Recv()

	err = ws.Close()
	Error(fmt.Sprintf("%T, %v\n", err, err))
	Error(err.(ConnectionCloseError).Msg)
}

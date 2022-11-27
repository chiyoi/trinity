package maru

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/atmt/message"
	"github.com/chiyoi/trinity/pkg/onebot"
	onebot_message "github.com/chiyoi/trinity/pkg/onebot/message"
	"github.com/chiyoi/trinity/pkg/trinity"
	"github.com/chiyoi/trinity/pkg/websocket"
)

func OnebotServer(rdb *redis.Client, chanFromAtmt <-chan onebot_message.Message) *http.Server {
	selectEvent := func(onebotEvent, atmtEvent <-chan onebot_message.Message) (data []byte, err error) {
		var msg onebot_message.Message
		select {
		case msg = <-onebotEvent:
		case msg = <-atmtEvent:
		}
		done := make(chan struct{})
		go func() {
			defer close(done)
			data, err = ws.Recv()
		}()

		select {
		case msg := <-chanFromAtmt:
			var ids []onebot.UserId
			ids, err = GetLoggedInList(rdb)
			if err != nil {
				return
			}
			for _, id := range ids {
				if err1 := onebot.SendMsg(ws, id, msg); err1 != nil {
					logs.Warning("maru:", err1)
					continue
				}
			}
			return
		case <-done:
			if err != nil {
				return
			}
			return
		}
	}
	serveWsConnection := func(ws websocket.WebSocket) error {
		defer func() { _ = ws.Close() }()
		for {
			data, err := selectEvent(ws)
			if err != nil {
				return err
			}
			if len(data) == 0 {
				continue
			}

			var ev onebot.Event
			if err = json.Unmarshal(data, &ev); err != nil {
				return err
			}
			if ev.PostType != onebot.EventMessage {
				logs.Info("maru: non-message event received.")
				continue
			}

			go serveEvent(ev)
		}
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Hijack(w, r)
		if err != nil {
			logs.Error("maru:", err)
			return
		}

		go func() {
			defer func() { _ = ws.Close() }()
			if err := serveWsConnection(ws); true {
				if cce, ok := err.(websocket.ConnectionCloseError); ok && cce.Code() == websocket.NormalClosure {
					logs.Info("maru: ws closed.")
					return
				}
				logs.Error(err)
				return
			}
		}()
	}
	return &http.Server{
		Addr:    ":http",
		Handler: http.HandlerFunc(handler),
	}
}

func EventServer(rdb *redis.Client) func(ev onebot.Event) {
	return func(ev onebot.Event) {
		var msg message.Message
		for _, m := range ev.Message {
			switch m.Type {
			case onebot_message.TypeText:
				msg.Append(message.Text(m.Data["text"]))
			case onebot_message.TypeImage:
				msg.Append(message.Image(path.Base(m.Data["file"]), m.Data["url"]))
			case onebot_message.TypeRecord:
				msg.Append(message.Record(path.Base(m.Data["file"]), m.Data["url"]))
			case onebot_message.TypeVideo:
				msg.Append(message.Video(path.Base(m.Data["file"]), m.Data["file"]))
			default:
				msg.Append(message.Text("unsupported segment"))
			}
		}

		ss := strings.Split(msg.Plaintext(), " ")
		if len(ss) == 3 && ss[0] == "login" {
			user, passwd := ss[1], ss[2]
			Login(rdb, ev.UserId, user, passwd)
			return
		}

		auth := GetAuthFromLoggedIn(rdb, ev.UserId)
		if auth == "" {
			return
		}

		if _, err := trinity.PostMessage(trinityUrl, auth, msg); err != nil {
			logs.Warning("maru:", err)
		}
	}
}

func OnebotWorker() func(ws websocket.WebSocket, onebotEvent chan<- onebot.Event) {
	return func(ws websocket.WebSocket, onebotEvent chan<- onebot.Event) {
		defer func() {
			if err := ws.Close(); err != nil {
				logs.Error("maru:", err)
			}
		}()
		for {
			data, err := ws.Recv()
			if err != nil {
				if cce, ok := err.(websocket.ConnectionCloseError); ok && cce.Code() == websocket.NormalClosure {
					logs.Info("maru: ws closed.")
				}
				logs.Error("maru:", err)
				return
			}

			var ev onebot.Event
			if err = json.Unmarshal(data, &ev); err != nil {
				logs.Error("maru:", err)
				continue
			}
			if ev.PostType != onebot.EventMessage {
				logs.Info("maru: non-message event received.")
				continue
			}
			onebotEvent <- ev
		}
	}
}

func AtmtServer(atmtEvent chan<- atmt.Event) *atmt.Server {
	handler := func(ev atmt.Event) {
		var msg onebot_message.Message
		msg.Append(onebot_message.Text(fmt.Sprintf("%s %s", ev.User, ev.Time.Format(time.RFC822))))
		for _, seg := range ev.Message {
			switch seg.Type {
			case message.TypeText:
				msg.Append(onebot_message.Text(seg.Data))
			case message.TypeImage:
				msg.Append(onebot_message.Image(seg.Ref.Url))
			case message.TypeRecord:
				msg.Append(onebot_message.Record(seg.Ref.Url))
			default:
				msg.Append(onebot_message.Text("unsupported message type"))
			}
		}
		atmtEvent <- atmt.Event{
			Time:      time.Time{},
			User:      "",
			MessageId: "",
			Message:   []message.Segment{},
		}
	}
	return &atmt.Server{
		Addr: ":8080",
		Handler: atmt.HandlerFunc(func(ev atmt.Event) {
			atmtEvent <- ev
		}),
	}
}

func StartOnebotSrv(srv *http.Server) {
	logs.Info("maru: listening", srv.Addr)
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		logs.Error(err)
		return
	}
	logs.Info(fmt.Sprintf("maru: server at %s closed.", srv.Addr))
}

func StopOnebotSrv(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logs.Error(err)
		return
	}
}

func StartAtmtSrv(srv *atmt.Server) {
	logs.Info("maru: listening", srv.Addr)
	// TODO: register self to trinity listeners
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		logs.Error(err)
		return
	}
	logs.Info(fmt.Sprintf("maru: server at %s closed.", srv.Addr))
}

func StopAtmtSrv(srv *atmt.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logs.Error(err)
		return
	}
}

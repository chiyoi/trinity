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
	atmt2 "github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/atmt/message"
	"github.com/chiyoi/trinity/pkg/onebot"
	onebot_message "github.com/chiyoi/trinity/pkg/onebot/message"
	"github.com/chiyoi/trinity/pkg/trinity"
	"github.com/chiyoi/trinity/pkg/websocket"
)

func OnebotServer(rdb *redis.Client, chanFromAtmt <-chan any) *http.Server {
	serveEvent := func(ev onebot.Event) {
		var msg message.Message
		for _, m := range ev.Message {
			var seg message.Segment
			switch m.Type {
			case onebot.MessageText:
				seg.Type = message.TypeText
				seg.Data = m.Data["text"]
			case onebot.MessageImage:
				seg.Type = message.TypeImage
				seg.Ref.Name = path.Base(m.Data["file"])
				seg.Ref.Url = m.Data["url"]
			case onebot.MessageRecord:
				seg.Type = message.TypeImage
				seg.Ref.Name = path.Base(m.Data["file"])
				seg.Ref.Url = m.Data["url"]
			case onebot.MessageVideo:
				seg.Type = message.TypeVideo
				seg.Ref.Name = path.Base(m.Data["file"])
				seg.Ref.Url = m.Data["file"]
			default:
				logs.Info("maru: unusual segment received.")
				continue
			}
			msg = append(msg, seg)
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
	pullEvent := func(ws websocket.WebSocket) (data []byte, err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("pull event: %w", err)
			}
		}()

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
			data, err := pullEvent(ws)
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

func AtmtServer(chanToOnebot chan<- any) *atmt2.Server {
	handler := func(ev atmt2.Event) {
		chanToOnebot <- fmt.Sprintf("%s %s", ev.User, ev.Time.Format(time.RFC822))
		for _, seg := range ev.Message {
			switch seg.Type {
			case message.TypeText:
				chanToOnebot <- onebot_message.Text(seg.Data)
			case message.TypeImage:
				chanToOnebot <- onebot_message.Image(seg.Ref.Url)
			case message.TypeRecord:
				chanToOnebot <- onebot_message.Record(seg.Ref.Url)
			default:
				chanToOnebot <- onebot_message.Text("unsupported message type")
			}
		}
	}
	return &atmt2.Server{
		Addr:    ":8080",
		Handler: atmt2.HandlerFunc(handler),
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

func StartAtmtSrv(srv *atmt2.Server) {
	logs.Info("maru: listening", srv.Addr)
	// TODO: register self to trinity listeners
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		logs.Error(err)
		return
	}
	logs.Info(fmt.Sprintf("maru: server at %s closed.", srv.Addr))
}

func StopAtmtSrv(srv *atmt2.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logs.Error(err)
		return
	}
}

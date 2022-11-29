package adapter

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
	"github.com/chiyoi/trinity/pkg/sdk/onebot"
	onebot_message "github.com/chiyoi/trinity/pkg/sdk/onebot/message"
	"github.com/chiyoi/trinity/pkg/sdk/trinity"
	"github.com/chiyoi/websocket"
)

func OnebotServer(rdb *redis.Client) func(ws websocket.WebSocket) {
	eventHandler := func(ev onebot.Event) {
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
			if err := Login(rdb, ev.UserId, user, passwd); err != nil {
				logs.Error("maru:", err)
				return
			}
			return
		}

		auth, err := GetAuthFromLoggedIn(rdb, ev.UserId)
		if err != nil {
			logs.Warning("maru:", err)
			return
		}

		if _, err := trinity.PostMessage(trinityUrl, auth, msg); err != nil {
			logs.Warning("maru:", err)
		}
	}

	return func(ws websocket.WebSocket) {
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
					return
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
			go eventHandler(ev)
		}
	}
}

func Server(rdb *redis.Client) *atmt.Server {
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

		ids, err := GetLoggedInList(rdb)
		if err != nil {
			logs.Error("maru:", err)
			return
		}

		for _, id := range ids {
			if err := OnebotSendMsg(id, msg); err != nil {
				logs.Error("maru:", err)
				return
			}
		}
	}
	return &atmt.Server{
		Addr:    ":8080",
		Handler: atmt.HandlerFunc(handler),
	}
}

func StartSrv(srv *atmt.Server) {
	logs.Info("maru: listening", srv.Addr)
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		logs.Error(err)
		return
	}
	logs.Info(fmt.Sprintf("maru: server at %s closed.", srv.Addr))
}

func StopSrv(srv *atmt.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logs.Error(err)
		return
	}
}

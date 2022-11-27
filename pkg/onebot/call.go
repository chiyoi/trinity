package onebot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chiyoi/trinity/pkg/onebot/message"
	"github.com/chiyoi/trinity/pkg/websocket"
)

func FormatMsg(a ...any) (msg message.Message) {
	for _, aa := range a {
		switch ta := aa.(type) {
		case message.Message:
			msg.Extend(ta)
		case message.Segment:
			msg.Append(ta)
		default:
			msg.Append(message.Text(fmt.Sprint(aa)))
		}
	}
	return
}

func SendMsgCtx(ctx context.Context, ws websocket.WebSocket, id UserId, a ...any) (err error) {
	req := Request{
		Action: ActionSendMsg,
		Params: map[string]any{
			"message.Message_type": MessagePrivate,
			"user_id":              id,
			"message.Message":      FormatMsg(a...),
		},
	}
	b, err := json.Marshal(req)
	if err != nil {
		return
	}
	return ws.SendCtx(ctx, websocket.OpTextFrame, b)
}

func SendMsg(ws websocket.WebSocket, id UserId, a ...any) (err error) {
	return SendMsgCtx(context.Background(), ws, id, a...)
}

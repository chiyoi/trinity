package request

import (
	"encoding/json"

	"github.com/chiyoi/trinity/internal/app/trinity/a11n"
	"github.com/chiyoi/trinity/internal/app/trinity/client"
	"github.com/chiyoi/trinity/internal/app/trinity/db"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/sdk/trinity"
)

func handlePostMessage(resp *atmt.Message, req atmt.DataRequest[trinity.Action], content []atmt.Paragraph) {
	logPrefix := "handle post message:"
	var args trinity.ArgsPostMessage
	if err := json.Unmarshal(req.Args, &args); err != nil {
		logs.Warning(logPrefix, "bad request.")
		atmt.Error(resp, atmt.StatusBadRequest)
		return
	}
	user, pass, err := a11n.VerifyAuthorization(args.Auth)
	if err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}
	if !pass {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusUnauthorized)
		return
	}
	id, err := db.PostMessage(db.Message{
		SenderId: user,
		Content:  content,
	})
	if err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}

	msg, err := (&atmt.MessageBuilder[atmt.DataPost]{
		Type: atmt.MessagePost,
		Data: atmt.DataPost{
			MessageId: id,
			SenderId:  user,
		},
		Content: content,
	}).Message()
	if err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}

	if err = client.PushMessageToListeners(msg); err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}
}

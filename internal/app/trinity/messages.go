package trinity

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/chiyoi/neko03/pkg/neko"
	"github.com/chiyoi/trinity/internal/app/trinity/client"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageId = primitive.ObjectID

func handlePostMessage(baseCtx context.Context, w http.ResponseWriter, req Request, coll *mongo.Collection, rdb *redis.Client, now time.Time, user string) {
	var reqData ReqDataPostMessage
	if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
		logs.Warning("cannot parse request data")
		neko.BadRequest(w)
		return
	}

	id, err := postMessage(baseCtx, coll, dMessage{
		Time:    now.Unix(),
		User:    user,
		Message: reqData.Message,
	})
	if err != nil {
		logs.Error(err)
		neko.InternalServerError(w)
		return
	}

	if err = client.PushEventToListeners(baseCtx, rdb, atmt.Event{
		Time:      now,
		User:      user,
		MessageId: id.Hex(),
		Message:   reqData.Message,
	}); err != nil {
		logs.Error(err)
		neko.InternalServerError(w)
		return
	}

	respData := RespDataPostMessage{
		MessageId: id.Hex(),
	}
	resp := Response{
		StatusCode: StatusOK,
		Data:       respData,
	}
	respBody, err := json.Marshal(resp)
	if err != nil {
		logs.Error(err)
		neko.InternalServerError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func handleGetMessage(baseCtx context.Context, w http.ResponseWriter, req Request, coll *mongo.Collection) {
	var reqData ReqDataGetMessage
	if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
		logs.Warning("cannot parse request data")
		neko.BadRequest(w)
		return
	}

	oid, err := primitive.ObjectIDFromHex(reqData.Id)
	if err != nil {
		logs.Error(err)
		neko.InternalServerError(w)
		return
	}

	doc, err := getMessage(baseCtx, coll, oid)
	if err != nil {
		logs.Error(err)
		neko.InternalServerError(w)
		return
	}

	respData := RespDataGetMessage{
		Time:      doc.Time,
		User:      doc.User,
		MessageId: doc.Id.Hex(),
		Message:   doc.Message,
	}
	respBody, err := json.Marshal(Response{
		StatusCode: StatusOK,
		Data:       respData,
	})
	if err != nil {
		logs.Error(err)
		neko.InternalServerError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func handleQueryMessageIdsTimeRange(baseCtx context.Context, w http.ResponseWriter, req Request, coll *mongo.Collection) {
	var reqData ReqDataQueryMessageTimeRange
	if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
		logs.Warning("cannot parse request data")
		neko.BadRequest(w)
		return
	}

	ids, err := queryMessageIds(baseCtx, coll, reqData.From, reqData.To)
	if err != nil {
		logs.Error(err)
		neko.InternalServerError(w)
		return
	}

	var hIds []string
	for _, id := range ids {
		hIds = append(hIds, id.Hex())
	}
	respData := RespDataQueryMessageTimeRange{
		Ids: hIds,
	}
	respBody, err := json.Marshal(Response{
		StatusCode: StatusOK,
		Data:       respData,
	})
	if err != nil {
		logs.Error(err)
		neko.InternalServerError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

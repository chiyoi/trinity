package trinity

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/chiyoi/trinity/internal/app/trinity/client"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

func handlePostMessage(baseCtx context.Context, w http.ResponseWriter, req Request, coll *mongo.Collection, rdb *redis.Client, now time.Time, user string) {
	var reqData ReqDataPostMessage
	if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
		badRequestCallback(w, err)
		return
	}

	unixNow := now.Unix()
	id, err := postMessage(baseCtx, coll, dMessage{
		Time:    &unixNow,
		User:    user,
		Message: reqData.Message,
	})
	if err != nil {
		internalServerErrorCallback(w, err)
		return
	}

	if err = client.PushEventToListeners(baseCtx, rdb, atmt.Event{
		Time:      now,
		User:      user,
		MessageId: id,
		Message:   reqData.Message,
	}); err != nil {
		internalServerErrorCallback(w, err)
		return
	}

	respData := RespDataPostMessage{
		MessageId: id,
	}
	resp := Response{
		StatusCode: StatusOK,
		Data:       respData,
	}
	respBody, err := json.Marshal(resp)
	if err != nil {
		internalServerErrorCallback(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func handleGetMessage(baseCtx context.Context, w http.ResponseWriter, req Request, coll *mongo.Collection) {
	var reqData ReqDataGetMessage
	if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
		badRequestCallback(w, err)
		return
	}

	doc, err := getMessage(baseCtx, coll, reqData.Id)
	if err != nil {
		internalServerErrorCallback(w, err)
		return
	}

	respData := RespDataGetMessage{
		Time:      *doc.Time,
		User:      doc.User,
		MessageId: doc.Id,
		Message:   doc.Message,
	}
	respBody, err := json.Marshal(Response{
		StatusCode: StatusOK,
		Data:       respData,
	})
	if err != nil {
		internalServerErrorCallback(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func handleQueryMessageIdsTimeRange(baseCtx context.Context, w http.ResponseWriter, req Request, coll *mongo.Collection) {
	var reqData ReqDataQueryMessageTimeRange
	if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
		badRequestCallback(w, err)
		return
	}

	ids, err := queryMessageIds(baseCtx, coll, reqData.From, reqData.To)
	if err != nil {
		internalServerErrorCallback(w, err)
		return
	}

	respData := RespDataQueryMessageTimeRange{
		Ids: ids,
	}
	respBody, err := json.Marshal(Response{
		StatusCode: StatusOK,
		Data:       respData,
	})
	if err != nil {
		internalServerErrorCallback(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

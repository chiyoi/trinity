package handlers

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"

// 	"github.com/chiyoi/neko03/pkg/neko"
// 	"github.com/chiyoi/trinity/internal/pkg/logs"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// func handleGetMessage(baseCtx context.Context, w http.ResponseWriter, req Request, coll *mongo.Collection) {
// 	var reqData ReqDataGetMessage
// 	if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
// 		logs.Warning("cannot parse request data")
// 		neko.BadRequest(w)
// 		return
// 	}

// 	oid, err := primitive.ObjectIDFromHex(reqData.Id)
// 	if err != nil {
// 		logs.Error(err)
// 		neko.InternalServerError(w)
// 		return
// 	}

// 	doc, err := getMessage(baseCtx, coll, oid)
// 	if err != nil {
// 		logs.Error(err)
// 		neko.InternalServerError(w)
// 		return
// 	}

// 	respData := RespDataGetMessage{
// 		Time:      doc.Time,
// 		User:      doc.User,
// 		MessageId: doc.Id.Hex(),
// 		Message:   doc.Message,
// 	}
// 	respBody, err := json.Marshal(Response{
// 		StatusCode: StatusOK,
// 		Data:       respData,
// 	})
// 	if err != nil {
// 		logs.Error(err)
// 		neko.InternalServerError(w)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(respBody)
// }

// func handleQueryMessageIdsTimeRange(baseCtx context.Context, w http.ResponseWriter, req Request, coll *mongo.Collection) {
// 	var reqData ReqDataQueryMessageTimeRange
// 	if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
// 		logs.Warning("cannot parse request data")
// 		neko.BadRequest(w)
// 		return
// 	}

// 	ids, err := queryMessageIds(baseCtx, coll, reqData.From, reqData.To)
// 	if err != nil {
// 		logs.Error(err)
// 		neko.InternalServerError(w)
// 		return
// 	}

// 	var hIds []string
// 	for _, id := range ids {
// 		hIds = append(hIds, id.Hex())
// 	}
// 	respData := RespDataQueryMessageTimeRange{
// 		Ids: hIds,
// 	}
// 	respBody, err := json.Marshal(Response{
// 		StatusCode: StatusOK,
// 		Data:       respData,
// 	})
// 	if err != nil {
// 		logs.Error(err)
// 		neko.InternalServerError(w)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(respBody)
// }

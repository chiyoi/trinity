package trinity

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/chiyoi/trinity/internal/app/trinity/config"
	"github.com/chiyoi/trinity/internal/pkg/logs"

	"github.com/chiyoi/neko03/pkg/neko"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Server(mongodb *mongo.Database, rdb *redis.Client) *http.Server {
	bg := context.Background()
	messageCollection, nekoCollection :=
		mongodb.Collection(config.Get[string]("MongodbCollectionMessages")),
		mongodb.Collection(nekos)

	handler := func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		user, ok := verifyAuth(bg, w, r, nekoCollection)
		if !ok {
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			logs.Error(err)
			neko.InternalServerError(w)
			ok = false
			return
		}
		var req Request
		if err = json.Unmarshal(data, &req); err != nil {
			logs.Error("cannot parse request:", err)
			neko.BadRequest(w)
			ok = false
			return
		}

		switch req.Action {
		case ActionPostMessage:
			handlePostMessage(bg, w, req, messageCollection, rdb, now, user)
		case ActionGetMessage:
			handleGetMessage(bg, w, req, messageCollection)
		case ActionQueryMessageIdsTimeRange:
			handleQueryMessageIdsTimeRange(bg, w, req, messageCollection)
		case ActionCacheFile:
			handleCacheFile(bg, w, req)
		case ActionVerifyAuthorization:
			respBody, err := json.Marshal(Response{
				StatusCode: StatusOK,
				Data:       RespDataVerifyAuthorization{},
			})
			if err != nil {
				logs.Error(err)
				neko.InternalServerError(w)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(respBody)
		default:
			logs.Warning("invalid action.")
			neko.BadRequest(w)
			return
		}
	}
	return &http.Server{
		Addr:    ":http",
		Handler: http.HandlerFunc(handler),
	}
}

func unauthorized(w http.ResponseWriter) {
	http.Error(w, "401 unauthorized", http.StatusUnauthorized)
}
func verifyAuth(baseCtx context.Context, w http.ResponseWriter, r *http.Request, coll *mongo.Collection) (user string, ok bool) {
	auth := strings.Split(r.Header.Get("Authorization"), " ")
	if len(auth) != 2 || strings.ToLower(auth[0]) != "basic" {
		logs.Warning("bad authorization header:", auth)
		r.Header.Set("WWW-Authorization", "Basic")
		unauthorized(w)
		ok = false
		return
	}
	b, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		logs.Warning("cannot decode token:", auth[1])
		r.Header.Set("WWW-Authorization", "Basic")
		unauthorized(w)
		ok = false
		return
	}
	t := strings.Split(string(b), ":")
	user, token := t[0], t[1]

	var res dUser
	ctx, cancel := context.WithTimeout(baseCtx, reqTimeout)
	defer cancel()
	err = coll.FindOne(ctx, bson.M{"name": user}).Decode(&res)
	if err != nil {
		ok = false
		if errors.Is(err, mongo.ErrNoDocuments) {
			logs.Warning("forbad unknown user:", user)
			neko.Forbidden(w)
			return
		}
		logs.Error(err)
		neko.InternalServerError(w)
		return
	}
	if res.Token != token {
		logs.Warning("forbad unmatched user-passwd:", user, token)
		neko.Forbidden(w)
		ok = false
		return
	}
	ok = true
	return
}

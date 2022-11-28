package trinity

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/chiyoi/trinity/internal/app/trinity/config"
	"github.com/chiyoi/trinity/internal/pkg/logs"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func internalServerErrorCallback(w http.ResponseWriter, err error) {
	logs.Error("trinity:", err)
	http.Error(w, "500 internal server error", http.StatusInternalServerError)
}

func badRequestCallback(w http.ResponseWriter, err error) {
	logs.Error("trinity:", err)
	http.Error(w, "400 bad request", http.StatusBadRequest)
}

func unauthorizedCallback(w http.ResponseWriter, r *http.Request) {
	logs.Warning("trinity: authorization error")
	r.Header.Set("WWW-Authorization", "Basic")
	http.Error(w, "401 unauthorized", http.StatusUnauthorized)
}

func forbiddenCallback(w http.ResponseWriter) {
	logs.Warning("trinity: forbade request")
	http.Error(w, "403 forbidden", http.StatusForbidden)
}

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

		req, ok := parseBody(w, r)
		if !ok {
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
		default:
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}
	}
	return &http.Server{
		Addr:    ":http",
		Handler: http.HandlerFunc(handler),
	}
}

func verifyAuth(baseCtx context.Context, w http.ResponseWriter, r *http.Request, coll *mongo.Collection) (user string, ok bool) {
	auth := strings.Split(r.Header.Get("Authorization"), " ")
	if len(auth) != 2 || strings.ToLower(auth[0]) != "basic" {
		unauthorizedCallback(w, r)
		ok = false
		return
	}
	b, err := base64.RawStdEncoding.DecodeString(auth[1])
	if err != nil {
		unauthorizedCallback(w, r)
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
			forbiddenCallback(w)
			return
		}
		internalServerErrorCallback(w, err)
		return
	}
	if res.Token != token {
		forbiddenCallback(w)
		ok = false
		return
	}
	ok = true
	return
}

func parseBody(w http.ResponseWriter, r *http.Request) (api Request, ok bool) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		internalServerErrorCallback(w, err)
		ok = false
		return
	}

	if err = json.Unmarshal(data, &api); err != nil {
		badRequestCallback(w, err)
		ok = false
		return
	}
	ok = true
	return
}

func StartSrv(srv *http.Server) {
	logs.Info("trinity: listening", srv.Addr)
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		logs.Error(err)
		return
	}
	logs.Info(fmt.Sprintf("trinity: server at %s closed.", srv.Addr))
}

func StopSrv(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logs.Error(err)
		return
	}
}

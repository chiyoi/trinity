package trinity

import (
	"testing"

	"github.com/chiyoi/neko03/pkg/neko"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/sdk/trinity"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	mongodb *mongo.Database
	rdb     *redis.Client
)

func init() {
	var err error
	if mongodb, err = OpenMongo(); err != nil {
		logs.Fatal(err)
	}
	if rdb, err = OpenRedis(); err != nil {
		logs.Fatal(err)
	}
}

func TestVerifyAuthorization(t *testing.T) {
	srv := Server(mongodb, rdb)
	go neko.StartSrv(srv, false)
	defer neko.StopSrv(srv)

	auth := trinity.CreateAuthorization("chiyoi", "Chiyoi@trinity1")
	pass, err := trinity.VerifyAuthorization("http://localhost/", auth)
	if err != nil {
		t.Fatal(err)
	}
	if pass {
		t.Fatal("not pass")
	}
}

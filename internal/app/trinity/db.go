package trinity

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/chiyoi/trinity/internal/app/trinity/config"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	reqTimeout = time.Second * 10
)

func OpenMongo() (db *mongo.Database, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("open mongo: %w", err)
		}
	}()
	bg := context.Background()
	mongodbUri, err := config.GetErr[url.URL]("MongodbURI")
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(bg, reqTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbUri.String()))
	if err != nil {
		return
	}

	mongodbDatabase, err := config.GetErr[string]("MongodbDatabase")
	if err != nil {
		return
	}
	db = client.Database(mongodbDatabase)
	return
}

func OpenRedis() (rdb *redis.Client, err error) {
	opt, err := config.GetErr[*redis.Options]("RedisOptions")
	if err != nil {
		return
	}
	rdb = redis.NewClient(opt)
	return
}

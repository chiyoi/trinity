package trinity

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbOperationTimeout = time.Second * 10
)

func OpenMongo() (db *mongo.Database, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("open mongo: %w", err)
		}
	}()
	bg := context.Background()
	mongodbUri, err := GetConfig[string]("MongodbURI")
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(bg, dbOperationTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbUri))
	if err != nil {
		return
	}

	mongodbDatabase, err := GetConfig[string]("MongodbDatabase")
	if err != nil {
		return
	}
	db = client.Database(mongodbDatabase)
	return
}

func OpenRedis() (rdb *redis.Client, err error) {
	opt, err := GetConfig[*redis.Options]("RedisOptions")
	if err != nil {
		return
	}
	rdb = redis.NewClient(opt)
	return
}

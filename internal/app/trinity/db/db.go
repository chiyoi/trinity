package db

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/chiyoi/trinity/internal/app/trinity/config"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbTimeout = time.Second * 10
)

var (
	errMongodbNotSet = fmt.Errorf("mongodb not set")
	errRdbNotSet     = fmt.Errorf("rdb not set")
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

	ctx, cancel := context.WithTimeout(bg, dbTimeout)
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

var pool struct {
	rdb     *redis.Client
	mongodb *mongo.Database
	mu      sync.RWMutex
}

func SetDB(r *redis.Client, m *mongo.Database) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	if r != nil {
		pool.rdb = r
	}
	if m != nil {
		pool.mongodb = m
	}
}
func GetDB() (r *redis.Client, m *mongo.Database) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()
	return pool.rdb, pool.mongodb
}

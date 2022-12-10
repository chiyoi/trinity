package db

import (
	"context"
	"fmt"
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
	mongodbURI := config.Get[string]("MongodbURI")
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbURI))
	if err != nil {
		return
	}

	mongodbDatabase := config.Get[string]("MongodbDatabase")
	db = client.Database(mongodbDatabase)
	return
}

func OpenRedis() (rdb *redis.Client) {
	opt := config.Get[*redis.Options]("RedisOptions")
	return redis.NewClient(opt)
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

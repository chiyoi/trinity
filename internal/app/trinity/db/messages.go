package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/chiyoi/trinity/third_party/gocq/adapter/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PostMessage(msg Message) (id primitive.ObjectID, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("post message: %w", err)
		}
	}()
	collName, err := config.GetErr[string]("MongodbCollectionMessages")
	if err != nil {
		return
	}
	_, mongodb := GetDB()
	if mongodb == nil {
		err = errors.New("mongodb not set")
		return
	}
	coll := mongodb.Collection(collName)
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	resp, err := coll.InsertOne(ctx, msg)
	if err != nil {
		return
	}
	id, ok := resp.InsertedID.(primitive.ObjectID)
	if !ok {
		err = errors.New("unexpected non-ObjectID message id")
		return
	}
	return
}

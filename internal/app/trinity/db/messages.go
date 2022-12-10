package db

import (
	"context"
	"fmt"

	"github.com/chiyoi/trinity/internal/app/trinity/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func collectionMessages() (coll *mongo.Collection, err error) {
	_, mongodb := GetDB()
	if mongodb == nil {
		err = errMongodbNotSet
		return
	}
	collName := config.Get[string]("CollectionMessages")
	coll = mongodb.Collection(collName)
	return
}

func PostMessage(msg Message) (id MessageID, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("db: %w", err)
		}
	}()
	coll, err := collectionMessages()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	resp, err := coll.InsertOne(ctx, msg)
	if err != nil {
		return
	}
	id, ok := resp.InsertedID.(primitive.ObjectID)
	if !ok {
		err = fmt.Errorf("unexpected non-ObjectID message id")
		return
	}
	return
}

func GetMessage(id MessageID) (msg Message, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("db: %w", err)
		}
	}()
	coll, err := collectionMessages()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	err = coll.FindOne(ctx, Message{
		ID: id,
	}).Decode(&msg)
	return
}

func QueryMessageIdsLatestCount(count int) (ids []MessageID, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("db: %w", err)
		}
	}()
	coll, err := collectionMessages()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	cur, err := coll.Find(ctx, bson.D{}, options.Find().
		SetProjection(bson.M{"_id": 1}).
		SetSort(bson.M{"_id": -1}).
		SetLimit(int64(count)),
	)
	if err != nil {
		return
	}

	var msgs []Message
	ctx, cancel = context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	if err = cur.All(ctx, &msgs); err != nil {
		return
	}
	for _, msg := range msgs {
		ids = append(ids, msg.ID)
	}
	return
}

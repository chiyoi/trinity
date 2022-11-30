package trinity

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/chiyoi/trinity/pkg/atmt/message"
)

type dUser struct {
	Name  string `bson:"name,omitempty"`
	Token string `bson:"token,omitempty"`
}

type dMessage struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Time    int64              `bson:"time"`
	User    string             `bson:"user,omitempty"`
	Message message.Message    `bson:"message,omitempty"`
}

func postMessage(baseCtx context.Context, msgColl *mongo.Collection, msg dMessage) (id MessageId, err error) {
	ctx, cancel := context.WithTimeout(baseCtx, reqTimeout)
	defer cancel()
	resp, err := msgColl.InsertOne(ctx, msg)
	if err != nil {
		return
	}
	id = resp.InsertedID.(primitive.ObjectID)
	return
}

func getMessage(baseCtx context.Context, coll *mongo.Collection, id MessageId) (doc dMessage, err error) {
	ctx, cancel := context.WithTimeout(baseCtx, reqTimeout)
	defer cancel()

	if err = coll.FindOne(ctx, bson.M{"_id": id}).Decode(&doc); err != nil {
		return
	}
	return
}

func queryMessageIds(baseCtx context.Context, coll *mongo.Collection, from, to int64) (ids []MessageId, err error) {
	var rang bson.A
	if from != 0 {
		rang = append(rang, bson.M{"$gte": from})
	}
	if to != 0 {
		rang = append(rang, bson.M{"$lte": to})
	}
	ctx, cancel := context.WithTimeout(baseCtx, reqTimeout)
	defer cancel()
	cur, err := coll.Find(ctx, bson.M{
		"time": bson.M{
			"$and": rang,
		},
	}, options.Find().SetProjection(
		bson.M{"_id": 1},
	))
	if err != nil {
		return
	}

	var msgs []dMessage
	ctx, cancel = context.WithTimeout(baseCtx, reqTimeout)
	defer cancel()
	if err = cur.All(ctx, &msgs); err != nil {
		return
	}
	for _, msg := range msgs {
		ids = append(ids, msg.Id)
	}
	return
}

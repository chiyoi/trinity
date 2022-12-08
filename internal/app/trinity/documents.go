package trinity

// import (
// 	"context"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func getMessage(baseCtx context.Context, coll *mongo.Collection, id MessageId) (doc dMessage, err error) {
// 	ctx, cancel := context.WithTimeout(baseCtx, reqTimeout)
// 	defer cancel()

// 	if err = coll.FindOne(ctx, bson.M{"_id": id}).Decode(&doc); err != nil {
// 		return
// 	}
// 	return
// }

// func queryMessageIds(baseCtx context.Context, coll *mongo.Collection, from, to int64) (ids []MessageId, err error) {
// 	var rang bson.A
// 	if from != 0 {
// 		rang = append(rang, bson.M{"$gte": from})
// 	}
// 	if to != 0 {
// 		rang = append(rang, bson.M{"$lte": to})
// 	}
// 	ctx, cancel := context.WithTimeout(baseCtx, reqTimeout)
// 	defer cancel()
// 	cur, err := coll.Find(ctx, bson.M{
// 		"time": bson.M{
// 			"$and": rang,
// 		},
// 	}, options.Find().SetProjection(
// 		bson.M{"_id": 1},
// 	))
// 	if err != nil {
// 		return
// 	}

// 	var msgs []dMessage
// 	ctx, cancel = context.WithTimeout(baseCtx, reqTimeout)
// 	defer cancel()
// 	if err = cur.All(ctx, &msgs); err != nil {
// 		return
// 	}
// 	for _, msg := range msgs {
// 		ids = append(ids, msg.Id)
// 	}
// 	return
// }

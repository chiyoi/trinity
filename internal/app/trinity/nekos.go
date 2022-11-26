package trinity

import (
	"context"
	"crypto/sha256"
	"encoding/base64"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/chiyoi/trinity/internal/pkg/logs"
)

func AddNeko(db *mongo.Database, user, passwd string) (err error) {
	if len(user) == 0 || len(passwd) == 0 {
		logs.Warning("trinity: empty neko~")
		return
	}

	coll := db.Collection(mongodbCollectionNekos)
	bg := context.Background()

	sum := sha256.Sum256([]byte(passwd))
	token := base64.StdEncoding.EncodeToString(sum[:])

	ctx, cancel := context.WithTimeout(bg, dbOperationTimeout)
	defer cancel()
	if _, err = coll.InsertOne(ctx, dUser{
		Name:  user,
		Token: token,
	}); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			logs.Warning("trinity: duplicated neko~")
			return
		}
		logs.Error("trinity: insert neko error:", err)
		return
	}
	return
}

func UpdateNeko(db *mongo.Database, user, passwd string) (err error) {
	coll := db.Collection(mongodbCollectionNekos)
	bg := context.Background()

	sum := sha256.Sum256([]byte(passwd))
	token := base64.StdEncoding.EncodeToString(sum[:])

	ctx, cancel := context.WithTimeout(bg, dbOperationTimeout)
	defer cancel()
	res, err := coll.UpdateOne(ctx, dUser{
		Name: user,
	}, bson.M{"$set": dUser{
		Token: token,
	}})
	if err != nil {
		logs.Error("trinity: update neko error:", err)
		return
	}
	if res.MatchedCount == 0 {
		logs.Warning("trinity: attempting to update non-exist neko.")
	}
	return
}

func RemoveNeko(db *mongo.Database, user string) (err error) {
	coll := db.Collection(mongodbCollectionNekos)
	bg := context.Background()

	ctx, cancel := context.WithTimeout(bg, dbOperationTimeout)
	defer cancel()
	res, err := coll.DeleteOne(ctx, dUser{
		Name: user,
	})
	if err != nil {
		logs.Error("trinity: delete neko error:", err)
		return
	}
	if res.DeletedCount == 0 {
		logs.Warning("trinity: attempting to update non-exist neko.")
	}
	return
}

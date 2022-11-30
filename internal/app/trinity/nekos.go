package trinity

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/chiyoi/trinity/internal/app/trinity/config"
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

var (
	nekos = config.Get[string]("MongodbCollectionNekos")
)

func AddNeko(db *mongo.Database, user, passwd string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("add neko: %w", err)
		}
	}()

	if len(user) == 0 || len(passwd) == 0 {
		logs.Warning("zero value neko~")
		return
	}

	coll := db.Collection(nekos)
	bg := context.Background()

	sum := sha256.Sum256([]byte(passwd))
	token := base64.StdEncoding.EncodeToString(sum[:])

	ctx, cancel := context.WithTimeout(bg, reqTimeout)
	defer cancel()
	if _, err = coll.InsertOne(ctx, dUser{
		Name:  user,
		Token: token,
	}); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			logs.Warning("duplicated neko~")
			return
		}
		return
	}
	return
}

func UpdateNeko(db *mongo.Database, user, passwd string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("update neko: %w", err)
		}
	}()

	coll := db.Collection(nekos)
	bg := context.Background()

	sum := sha256.Sum256([]byte(passwd))
	token := base64.StdEncoding.EncodeToString(sum[:])

	ctx, cancel := context.WithTimeout(bg, reqTimeout)
	defer cancel()
	res, err := coll.UpdateOne(ctx, dUser{
		Name: user,
	}, bson.M{"$set": dUser{
		Token: token,
	}})
	if err != nil {
		return
	}
	if res.MatchedCount == 0 {
		logs.Warning("neko not found~")
	}
	return
}

func RemoveNeko(db *mongo.Database, user string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("remove neko: %w", err)
		}
	}()

	coll := db.Collection(nekos)
	bg := context.Background()

	ctx, cancel := context.WithTimeout(bg, reqTimeout)
	defer cancel()
	res, err := coll.DeleteOne(ctx, dUser{
		Name: user,
	})
	if err != nil {
		return
	}
	if res.DeletedCount == 0 {
		logs.Warning("neko not found~")
	}
	return
}

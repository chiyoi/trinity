package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/chiyoi/trinity/internal/app/aira/config"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func VerifyUserToken(user string, token string) (pass bool, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("verify user token: %w", err)
		}
	}()
	collName, err := config.GetErr[string]("MongodbCollectionNekos")
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
	var u User
	if err = coll.FindOne(ctx, bson.M{"name": user}).Decode(&u); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logs.Warning("forbad unknown user:", user)
			pass = false
			err = nil
			return
		}
		return
	}
	if u.Token != token {
		logs.Warning("forbad unmatched user-passwd:", user, token)
		pass = false
		return
	}
	pass = true
	return
}

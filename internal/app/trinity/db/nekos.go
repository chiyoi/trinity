package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/chiyoi/trinity/internal/app/trinity/config"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/sdk/trinity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func collectionNekos() (coll *mongo.Collection, err error) {
	collName, err := config.GetErr[string]("MongodbCollectionNekos")
	if err != nil {
		return
	}
	_, mongodb := GetDB()
	if mongodb == nil {
		err = errMongodbNotSet
		return
	}
	coll = mongodb.Collection(collName)
	return
}

func VerifyUserToken(user string, token string) (pass bool, err error) {
	logPrefix := "verify user token:"
	defer func() {
		if err != nil {
			err = fmt.Errorf("verify user token: %w", err)
		}
	}()
	coll, err := collectionNekos()
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var u User
	if err = coll.FindOne(ctx, bson.M{"name": user}).Decode(&u); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logs.Warning(logPrefix, "forbad unknown user:", user)
			pass = false
			err = nil
			return
		}
		return
	}
	if u.Token != token {
		logs.Warning(logPrefix, "forbad unmatched user-passwd:", user, token)
		pass = false
		return
	}
	pass = true
	return
}

func AddNeko(user, passwd string) (err error) {
	logPrefix := "add neko:"
	defer func() {
		if err != nil {
			err = fmt.Errorf("add neko: %w", err)
		}
	}()
	if len(user) == 0 || len(passwd) == 0 {
		logs.Warning(logPrefix, "zero value neko~")
		return
	}
	coll, err := collectionNekos()
	if err != nil {
		return
	}

	token := trinity.PasswdToken(passwd)
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	if _, err = coll.InsertOne(ctx, User{
		Name:  user,
		Token: token,
	}); err != nil {
		return
	}
	return
}

func UpdateNeko(user, passwd string) (err error) {
	logPrefix := "update neko:"
	defer func() {
		if err != nil {
			err = fmt.Errorf("update neko: %w", err)
		}
	}()
	coll, err := collectionNekos()
	if err != nil {
		return
	}

	token := trinity.PasswdToken(passwd)
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	res, err := coll.UpdateOne(ctx, User{
		Name: user,
	}, bson.M{"$set": User{
		Token: token,
	}})
	if err != nil {
		return
	}
	if res.MatchedCount == 0 {
		logs.Warning(logPrefix, "neko not found~")
	}
	return
}

func RemoveNeko(user string) (err error) {
	logPrefix := "remove neko:"
	defer func() {
		if err != nil {
			err = fmt.Errorf("remove neko: %w", err)
		}
	}()
	coll, err := collectionNekos()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	res, err := coll.DeleteOne(ctx, User{
		Name: user,
	})
	if err != nil {
		return
	}
	if res.DeletedCount == 0 {
		logs.Warning(logPrefix, "neko not found~")
	}
	return
}

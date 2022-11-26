package trinity

import (
	"github.com/chiyoi/trinity"
	"github.com/chiyoi/trinity/internal/configs"
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

var cfg = trinity.TrinityConfig

func GetConfig[T any](key string) (a T, err error) {
	return configs.Get[T](cfg, key)
}

var (
	mongodbCollectionNekos, mongodbCollectionMessages string
)

func init() {
	var err error
	mongodbCollectionNekos, err = GetConfig[string]("MongodbCollectionNekos")
	if err != nil {
		logs.Fatal("trinity:", err)
	}
	mongodbCollectionMessages, err = GetConfig[string]("MongodbCollectionMessages")
	if err != nil {
		logs.Fatal("trinity:", err)
	}

}

package trinity

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/chiyoi/trinity/pkg/trinity"
	"github.com/go-redis/redis/v8"
)

var TrinityConfig = map[string]any{
	"ServiceURL": "http://trinity/",

	"MongodbURI":                "to be init",
	"MongodbDatabase":           "trinity",
	"MongodbCollectionNekos":    "nekos",
	"MongodbCollectionMessages": "messages",

	"RedisOptions": &redis.Options{
		Addr:     "redis-18080.c56.east-us.azure.cloud.redislabs.com:18080",
		Username: "trinity",
		Password: "Neko03Trinity@redis",
	},
	"RedisKeyListeners": "trinity:listeners",

	"AzureStorageAccount": "neko03storage",
	"AzureStorageKey":     "lZzvHnmRwYiD1t9xEDZhxn07eNtmn4J3qiu/8UGkfGEeL1Pz3C/yR8+hY7rmJo/xVuTLMtilsq/7+ASte3hwBQ==",
	"FileCacheContainer":  "trinity-file-cache",
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	keyFile := filepath.Join(home, ".ssh/keys/chiyoi_atlas")

	q := url.Values{}
	q.Set("authSource", "$external")
	q.Set("authMechanism", "MONGODB-X509")
	q.Set("retryWrites", "true")
	q.Set("w", "majority")
	q.Set("tlsCertificateKeyFile", keyFile)
	TrinityConfig["MongodbURI"] = (&url.URL{
		Scheme:   "mongodb+srv",
		Host:     "cluster0.catoops.mongodb.net",
		Path:     "/",
		RawQuery: q.Encode(),
	}).String()
}

var AiraConfig = map[string]any{
	"ServiceURL": "http://aira/",

	"TrinityURL": "http://trinity/",
	"TrinityAccessToken": trinity.CreateAuthorization(
		"aira",
		"Neko03Aira@trinity",
	),

	"RedisOptions": &redis.Options{
		Addr:     "redis-18080.c56.east-us.azure.cloud.redislabs.com:18080",
		Username: "aira",
		Password: "Neko03Aira@redis",
	},
	"RedisKeyListeners": "trinity:listeners",

	"CommandPrefix": []string{"."},
}

var MaruConfig = map[string]any{
	"ServiceURL": "http://maru/",

	"TrinityURL":     "http://trinity/",
	"OnebotURL":      "http://gocq/",
	"OnebotEventURL": "http://gocq:8080/",

	"RedisOptions": &redis.Options{
		Addr:     "redis-18080.c56.east-us.azure.cloud.redislabs.com:18080",
		Username: "maru",
		Password: "Neko03Maru@redis",
	},
	"RedisKeyUsersLoggedIn": "maru:usersLoggedIn",
	"RedisKeyNekoMap":       "maru:nekoMap",
	"RedisKeyListeners":     "trinity:listeners",
}

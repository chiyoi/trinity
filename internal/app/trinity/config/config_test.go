package config

import "testing"

func TestGetErr(t *testing.T) {
	nekos, err := GetErr[string]("MongodbCollectionNekos")
	if err != nil {
		t.Error(err)
	}
	if nekos != "nekos" {
		t.Error(nekos)
	}
}

package db

import "testing"

func TestCollectionNekos(t *testing.T) {
	nekos, err := collectionNekos()
	if err != nil {
		t.Error(err)
	}
	t.Log(nekos)
}

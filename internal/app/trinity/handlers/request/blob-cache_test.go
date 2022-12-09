package request

import (
	"os"
	"testing"

	"github.com/chiyoi/trinity/pkg/sdk/trinity"
)

func TestGetBlobCacheURL(t *testing.T) {
	auth := trinity.CreateAuthorization("chiyoi", "Chiyoi@trinity")
	vals, _, err := trinity.Request[trinity.ArgsGetBlobCacheURL, trinity.ValuesGetBlobCacheURL](
		"http://localhost/",
		trinity.ActionGetBlobCacheURL,
		trinity.ArgsGetBlobCacheURL{
			Auth:     auth,
			BlobName: "nyan",
		},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vals.SasURL)
}

func TestCacheBlob(t *testing.T) {
	auth := trinity.CreateAuthorization("chiyoi", "Chiyoi@trinity")
	b, err := os.ReadFile("/Users/chiyoi/Desktop/IMG_4618.png")
	if err != nil {
		t.Fatal(err)
	}
	sasURL, err := trinity.CacheBlob("http://localhost/", auth, b)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sasURL)
}

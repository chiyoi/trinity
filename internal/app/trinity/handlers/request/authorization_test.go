package request

import (
	"testing"

	"github.com/chiyoi/trinity/pkg/sdk/trinity"
)

func TestVerifyAuthorization(t *testing.T) {
	auth := trinity.CreateAuthorization("chiyoi", "Chiyoi@trinity")
	resp, _, err := trinity.Request[trinity.ArgsVerifyAuthorization, trinity.ValuesVerifyAuthorization](
		"http://localhost/",
		trinity.ActionVerifyAuthorization,
		trinity.ArgsVerifyAuthorization{
			Auth: auth,
		},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Pass {
		t.Fatal("not passed.")
	}
}

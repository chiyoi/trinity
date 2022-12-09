package request

import (
	"fmt"
	"testing"

	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/sdk/trinity"
)

func TestPostMessage(t *testing.T) {
	auth := trinity.CreateAuthorization("chiyoi", "Chiyoi@trinity")
	if _, _, err := trinity.Request[trinity.ArgsPostMessage, trinity.ValuesPostMessage](
		"http://localhost:3333/",
		trinity.ActionPostMessage,
		trinity.ArgsPostMessage{
			Auth: auth,
		},
		atmt.FormatContent("aira"),
	); err != nil {
		t.Fatal(err)
	}
}

func TestQueryMessageIdsLatestCount(t *testing.T) {
	auth := trinity.CreateAuthorization("chiyoi", "Chiyoi@trinity")
	vals, _, err := trinity.Request[trinity.ArgsQueryMessageIdsLatestCount, trinity.ValuesQueryMessageIds](
		"http://localhost/",
		trinity.ActionQueryMessageIdsLatestCount,
		trinity.ArgsQueryMessageIdsLatestCount{
			Auth:  auth,
			Count: 5,
		},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vals.Ids)
}

func TestGetMessage(t *testing.T) {
	auth := trinity.CreateAuthorization("chiyoi", "Chiyoi@trinity")
	vals1, _, err := trinity.Request[trinity.ArgsQueryMessageIdsLatestCount, trinity.ValuesQueryMessageIds](
		"http://localhost/",
		trinity.ActionQueryMessageIdsLatestCount,
		trinity.ArgsQueryMessageIdsLatestCount{
			Auth:  auth,
			Count: 1,
		},
		nil,
	)
	if err != nil || len(vals1.Ids) != 1 {
		if err == nil {
			err = fmt.Errorf("error length: %v", len(vals1.Ids))
		}
		t.Fatal(err)
	}
	vals2, ret, err := trinity.Request[trinity.ArgsGetMessage, trinity.ValuesGetMessage](
		"http://localhost/",
		trinity.ActionGetMessage,
		trinity.ArgsGetMessage{
			Auth: auth,
			ID:   vals1.Ids[0],
		},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vals2.MessageID.Timestamp(), vals2.Sender, ret)
}

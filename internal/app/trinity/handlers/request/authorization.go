package request

import (
	"encoding/json"

	"github.com/chiyoi/trinity/internal/app/trinity/db"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/sdk/trinity"
)

func verifyAuth(auth string) (user string, pass bool, err error) {
	user, passwdS256, err := trinity.ParseAuthorization(auth)
	if err != nil {
		return
	}
	token := trinity.PasswdS256Token(passwdS256)
	pass, err = db.VerifyUserToken(user, token)
	return
}

func handleVerifyAuthorization(resp *atmt.Message, req atmt.DataRequest[trinity.Action]) {
	logPrefix := "handle verify authorization:"
	var args trinity.ArgsVerifyAuthorization
	if err := json.Unmarshal(req.Args, &args); err != nil {
		logs.Warning(logPrefix, err)
		atmt.Error(resp, atmt.StatusBadRequest)
		return
	}
	_, pass, err := verifyAuth(args.Auth)
	if err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}

	b := atmt.MessageBuilder[atmt.DataResponseBuilder[trinity.ValuesVerifyAuthorization]]{
		Type: atmt.MessageResponse,
		Data: atmt.DataResponseBuilder[trinity.ValuesVerifyAuthorization]{
			StatusCode: atmt.StatusOK,
			Values: trinity.ValuesVerifyAuthorization{
				Pass: pass,
			},
		},
	}
	if err = b.Write(resp); err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}
}

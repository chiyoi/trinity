package request

import (
	"encoding/json"

	"github.com/chiyoi/neko03/pkg/logs"
	"github.com/chiyoi/trinity/internal/app/trinity/db"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/sdk/trinity"
)

func verifyAuth(resp *atmt.Message, auth string) (user string, pass bool, err error) {
	user, passwdS256, err := trinity.ParseAuthorization(auth)
	if err != nil {
		logs.Warningf("cannot parse authorization(%s)", err)
		atmt.BadRequest(resp)
		return
	}
	token := trinity.PasswdS256Token(passwdS256)
	if pass, err = db.VerifyUserToken(user, token); err != nil {
		logs.Error(err)
		atmt.InternalServerError(resp)
	}
	return
}

func handleVerifyAuthorization(resp *atmt.Message, req atmt.DataRequest[trinity.Action]) {
	var args trinity.ArgsVerifyAuthorization
	if err := json.Unmarshal(req.Args, &args); err != nil {
		logs.Warning(err)
		atmt.Error(resp, atmt.StatusBadRequest)
		return
	}
	_, pass, err := verifyAuth(resp, args.Auth)
	if err != nil {
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
		logs.Error(err)
		atmt.InternalServerError(resp)
		return
	}
}

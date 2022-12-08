package a11n

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/chiyoi/trinity/internal/app/trinity/db"
)

func VerifyAuthorization(auth string) (user string, pass bool, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("verify authorization: %w", err)
		}
	}()
	b, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return
	}
	t := strings.Split(string(b), ":")
	user, token := t[0], t[1]
	pass, err = db.VerifyUserToken(user, token)
	return
}

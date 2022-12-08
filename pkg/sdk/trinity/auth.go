package trinity

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

func CreateAuthorization(user, passwd string) (auth string) {
	passwdS256 := sha256.Sum256([]byte(passwd))
	passwdS256B64 := base64.StdEncoding.EncodeToString(passwdS256[:])
	token := fmt.Sprintf("%s:%s", user, passwdS256B64)
	return base64.StdEncoding.EncodeToString([]byte(token))
}

func ParseAuthToken(auth string) (user string, passwdS256 [sha256.Size]byte, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parse authorization token: %w", err)
		}
	}()
	token, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return
	}
	t := strings.Split(string(token), ":")
	if len(t) != 2 {
		err = errors.New("bad format")
	}
	user, passwdS256B64 := t[0], t[1]
	passwdS256Slice, err := base64.StdEncoding.DecodeString(passwdS256B64)
	if err != nil {
		return
	}
	passwdS256 = *(*[sha256.Size]byte)(passwdS256Slice)
	return
}

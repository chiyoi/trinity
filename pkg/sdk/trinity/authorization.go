package trinity

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

func PasswdToken(passwd string) (token string) {
	passwdS256 := sha256.Sum256([]byte(passwd))
	return PasswdS256Token(passwdS256)
}
func PasswdS256Token(passwdS256 [sha256.Size]byte) (token string) {
	return base64.StdEncoding.EncodeToString(passwdS256[:])
}

func CreateAuthorization(user, passwd string) (auth string) {
	token := PasswdToken(passwd)
	return CreateAuthorizationToken(user, token)
}
func CreateAuthorizationToken(user, token string) (auth string) {
	pair := fmt.Sprintf("%s:%s", user, token)
	return base64.StdEncoding.EncodeToString([]byte(pair))
}

func ParseAuthorization(auth string) (user string, passwdS256 [sha256.Size]byte, err error) {
	token, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return
	}
	t := strings.Split(string(token), ":")
	if len(t) != 2 {
		err = fmt.Errorf("bad format")
		return
	}
	user, passwdS256B64 := t[0], t[1]
	passwdS256Slice, err := base64.StdEncoding.DecodeString(passwdS256B64)
	if err != nil {
		return
	}
	passwdS256 = *(*[sha256.Size]byte)(passwdS256Slice)
	return
}

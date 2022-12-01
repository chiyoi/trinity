package trinity

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func CreateAuthorization(user, passwd string) string {
	passwdS256 := sha256.Sum256([]byte(passwd))
	passwdS256B64 := base64.StdEncoding.EncodeToString(passwdS256[:])
	token := fmt.Sprintf("%s:%s", user, passwdS256B64)
	tokenB64 := base64.StdEncoding.EncodeToString([]byte(token))
	return fmt.Sprintf("Basic %s", tokenB64)
}

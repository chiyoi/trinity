package trinity

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

func CreateAuthorization(user, passwd string) string {
	sum := sha256.Sum256([]byte(passwd))
	return strings.Join([]string{"Basic", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(
		"%s:%s",
		user,
		base64.StdEncoding.EncodeToString(sum[:]),
	)))}, " ")
}

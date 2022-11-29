package main

import (
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/trinity"
)

func main() {
	auth := trinity.CreateAuthorization("chiyoi", "Chiyoi@trinity")
	mid, err := trinity.PostMessage("http://localhost:3333/", auth, "hello, world!")
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Info(mid)
}

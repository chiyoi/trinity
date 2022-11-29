package main

import (
	"github.com/chiyoi/trinity/internal/app/trinity"
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

func main() {
	db, err := trinity.OpenMongo()
	if err != nil {
		logs.Fatal("update neko:", err)
	}
	trinity.UpdateNeko(db, "chiyoi", "Chiyoi@trinity")
	logs.Info("completed")
}

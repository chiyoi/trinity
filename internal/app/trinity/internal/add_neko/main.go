package main

import (
	"github.com/chiyoi/trinity/internal/app/trinity"
	"github.com/chiyoi/trinity/internal/pkg/logs"
)

func main() {
	db, err := trinity.OpenMongo()
	if err != nil {
		logs.Fatal("add neko:", err)
	}
	trinity.AddNeko(db, "chiyoi", "Chiyoi@Neko03Trinity")
	logs.Info("completed")
}

// chiyoi:Chiyoi@Neko03Trinity
// aira:Neko03Aira

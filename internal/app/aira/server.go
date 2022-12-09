package aira

import (
	"github.com/chiyoi/trinity/internal/app/aira/handlers"
	"github.com/chiyoi/trinity/internal/app/aira/handlers/aira"
	"github.com/chiyoi/trinity/pkg/atmt"
)

func Server() *atmt.Server {
	mux := atmt.NewServeMux()
	mux.Handle(aira.Aira())
	return &atmt.Server{
		Addr:    ":http",
		Handler: handlers.LogMessage(mux),
	}
}

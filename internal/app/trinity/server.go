package trinity

import (
	"github.com/chiyoi/trinity/internal/app/trinity/handlers"
	"github.com/chiyoi/trinity/internal/app/trinity/handlers/request"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/atmt/rules"
)

func Server() *atmt.Server {
	mux := atmt.NewServeMux()
	mux.Handle(rules.MessageType(atmt.MessageRequest).Matcher(), request.Handler())
	return &atmt.Server{
		Addr:    ":http",
		Handler: handlers.LogMessage(mux),
	}
}

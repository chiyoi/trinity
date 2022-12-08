package trinity

import (
	"github.com/chiyoi/trinity/internal/app/trinity/handlers/request"
	"github.com/chiyoi/trinity/pkg/atmt"
)

func Server() *atmt.Server {
	mux := atmt.NewServeMux()
	mux.Handle(request.Request())
	return &atmt.Server{
		Addr:    ":http",
		Handler: mux,
	}
}

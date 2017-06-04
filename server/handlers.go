package server

import (
	"net/http"

	"github.com/m3db/m3ctl/handler"
	"github.com/m3db/m3ctl/handler/rules"
)

var (
	healthURL   = "/health"
	controllers = []handler.Controller{
		rules.Controller{},
	}
)

func registerHandlers(mux *http.ServeMux) error {
	mux.HandleFunc(healthURL, handler.M3CtlHealthCheck)
	for _, controller := range controllers {
		controller.RegisterHandlers(mux)
	}
	return nil
}

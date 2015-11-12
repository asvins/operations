package main

import (
	"net/http"
	"strings"

	"github.com/asvins/router"
	"github.com/asvins/router/logger"
	"github.com/unrolled/render"
)

func DiscoveryHandler(w http.ResponseWriter, req *http.Request) {
	prefix := strings.Join([]string{ServerConfig.Server.Addr, ServerConfig.Server.Port}, ":")
	r := render.New()

	//add discovery links here
	discoveryMap := map[string]string{"discovery": prefix + "/api/discovery"}

	r.JSON(w, http.StatusOK, discoveryMap)
}

func DefRoutes() *router.Router {
	r := router.NewRouter()

	r.AddRoute("/api/discovery", router.GET, DiscoveryHandler)

	// interceptors
	r.AddBaseInterceptor("/", logger.NewLogger())
	return r
}

package main

import (
	"net/http"
	"strings"

	"github.com/asvins/router"
	"github.com/asvins/router/errors"
	"github.com/asvins/router/logger"
	"github.com/unrolled/render"
)

func DiscoveryHandler(w http.ResponseWriter, req *http.Request) errors.Http {
	prefix := strings.Join([]string{ServerConfig.Server.Addr, ServerConfig.Server.Port}, ":")
	r := render.New()

	//add discovery links here
	discoveryMap := map[string]string{"discovery": prefix + "/api/discovery",
		"active_packs": prefix + "/api/packs/active",
		"packs":        prefix + "/api/packs/all"}

	r.JSON(w, http.StatusOK, discoveryMap)
	return nil
}

func DefRoutes() *router.Router {
	r := router.NewRouter()

	r.Handle("/api/discovery", router.GET, DiscoveryHandler, []router.Interceptor{})
	r.Handle("/api/box", router.GET, retrieveBoxes, []router.Interceptor{})
	r.Handle("/api/packs", router.GET, retrievePacks, []router.Interceptor{})
	r.Handle("/api/packMedications", router.GET, retrievePackMedications, []router.Interceptor{})

	// interceptors
	r.AddBaseInterceptor("/", logger.NewLogger())
	return r
}

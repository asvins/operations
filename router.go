package main

import (
	"net/http"
	"strings"

	"github.com/asvins/common_db/postgres"
	"github.com/asvins/router"
	"github.com/asvins/router/logger"
	"github.com/unrolled/render"
)

func DiscoveryHandler(w http.ResponseWriter, req *http.Request) {
	prefix := strings.Join([]string{ServerConfig.Server.Addr, ServerConfig.Server.Port}, ":")
	r := render.New()

	//add discovery links here
	discoveryMap := map[string]string{"discovery": prefix + "/api/discovery",
		"active_packs": prefix + "/api/packs/active",
		"packs":        prefix + "/api/packs/all"}

	r.JSON(w, http.StatusOK, discoveryMap)
}

func ActivePacksHandler(w http.ResponseWriter, req *http.Request) {
	r := render.New()
	if req.ParseForm() != nil {
		http.Error(w, "Invalid Input", 400)
		return
	}

	owner := req.Form.Get("id")
	var ps []Pack
	db := postgres.GetDatabase(DatabaseConfig)
	if GetActivePacks(owner, &ps, db) != nil {
		http.NotFound(w, req)
	}

	r.JSON(w, http.StatusOK, ps)
}

func AllPacksHandler(w http.ResponseWriter, req *http.Request) {
	r := render.New()

	var ps []Pack
	if req.ParseForm() != nil {
		http.Error(w, "Invalid Input", 400)
		return
	}

	status := req.Form.Get("status")
	db := postgres.GetDatabase(DatabaseConfig)
	if GetPacksByStatusString(status, &ps, db) != nil {
		http.NotFound(w, req)
	}

	r.JSON(w, http.StatusOK, ps)
}

func DefRoutes() *router.Router {
	r := router.NewRouter()

	r.AddRoute("/api/discovery", router.GET, DiscoveryHandler)
	r.AddRoute("/api/packs/active", router.GET, ActivePacksHandler)
	r.AddRoute("/api/packs/all", router.GET, AllPacksHandler)

	// interceptors
	r.AddBaseInterceptor("/", logger.NewLogger())
	//r.AddBaseInterceptor("/api/packs", auth.HasOneScopeInterceptor([]string["admin", "patient"]))
	return r
}

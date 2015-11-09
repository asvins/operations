package main

import (
	"fmt"
	"net/http"

	"log"

	// "github.com/asvins/common_db/postgres"
	"github.com/asvins/utils/config"
)

var ServerConfig *Config = new(Config)
// var DatabaseConfig *postgres.Config

// function that will run before main
func init() {
	fmt.Println("[INFO] Initializing server")
	err := config.Load("operations_config.gcfg", ServerConfig)
	if err != nil {
		log.Fatal(err)
	}

	// DatabaseConfig = postgres.NewConfig(ServerConfig.Database.User, ServerConfig.Database.DbName, ServerConfig.Database.SSLMode)
	fmt.Println("[INFO] Initialization Done!")
}

func main() {
	fmt.Println("[INFO] Server running on port:", ServerConfig.Server.Port)
	r := DefRoutes()
	http.ListenAndServe(":"+ServerConfig.Server.Port, r)
}

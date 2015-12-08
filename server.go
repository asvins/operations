package main

import (
	"fmt"
	"net/http"

	"log"

	"github.com/asvins/common_db/postgres"
	"github.com/asvins/common_io"
	"github.com/asvins/utils/config"
)

var (
	ServerConfig   *Config = new(Config)
	DatabaseConfig *postgres.Config
	producer       *common_io.Producer
	consumer       *common_io.Consumer
)

func init() {
	fmt.Println("[INFO] Initializing server")
	err := config.Load("operations_config.gcfg", ServerConfig)
	if err != nil {
		log.Fatal(err)
	}

	/*
	*	Database
	 */
	DatabaseConfig = postgres.NewConfig(ServerConfig.Database.User, ServerConfig.Database.DbName, ServerConfig.Database.SSLMode)

	/*
	*	Common io
	 */
	setupCommonIo()

	fmt.Println("[INFO] Initialization Done!")
}

func main() {
	fmt.Println("[INFO] Server running on port:", ServerConfig.Server.Port)
	r := DefRoutes()
	http.ListenAndServe(":"+ServerConfig.Server.Port, r)
}

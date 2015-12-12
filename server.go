package main

import (
	"fmt"
	"net/http"

	"log"

	"github.com/asvins/common_db/postgres"
	"github.com/asvins/common_io"
	"github.com/asvins/utils/config"
	"github.com/jinzhu/gorm"
	"github.com/unrolled/render"
)

var (
	ServerConfig   *Config = new(Config)
	DatabaseConfig *postgres.Config
	db             *gorm.DB
	producer       *common_io.Producer
	consumer       *common_io.Consumer
	rend           *render.Render = render.New()
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
	db = postgres.GetDatabase(DatabaseConfig)

	/*
	*	Common io
	 */
	setupCommonIo()

	/*
	*	Crons
	 */
	startShipCron()
	startPackNotificationCron()

	fmt.Println("[INFO] Initialization Done!")
}

func main() {
	fmt.Println("[INFO] Server running on port:", ServerConfig.Server.Port)
	r := DefRoutes()
	http.ListenAndServe(":"+ServerConfig.Server.Port, r)
}

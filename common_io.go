package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/asvins/common_db/postgres"
	"github.com/asvins/common_io"
	"github.com/asvins/utils/config"
)

func setupCommonIo() {
	cfg := common_io.Config{}

	err := config.Load("common_io_config.gcfg", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	/*
	*	Producer
	 */
	producer, err = common_io.NewProducer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	/*
	*	Consumer
	 */
	consumer = common_io.NewConsumer(cfg)

	/*
	*	topics
	 */
	consumer.HandleTopic("treatment_scheduled", treatmentScheduledHandler)

	if err = consumer.StartListening(); err != nil {
		log.Fatal(err)
	}
}

func treatmentScheduledHandler(msg []byte) {
	p := Pack{}
	err := json.Unmarshal(msg, &p)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return
	}

	db := postgres.GetDatabase(DatabaseConfig)
	p.Create(db)

	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
	}
	producer.Publish("pack_created", b)
}

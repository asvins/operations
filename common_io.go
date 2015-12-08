package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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
	consumer.HandleTopic("treatment_created", treatmentCreatedHandler)

	if err = consumer.StartListening(); err != nil {
		log.Fatal(err)
	}
}

func treatmentCreatedHandler(msg []byte) {
	p := Pack{}
	err := json.Unmarshal(msg, &p)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return
	}

	p.From = time.Now()
	p.To = time.Now().AddDate(0, 1, 0)

	db := postgres.GetDatabase(DatabaseConfig)
	p.Create(db)

	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
	}
	producer.Publish("pack_created", b)
}

package main

import (
	"log"

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
	consumer.HandleTopic("subscription_paid", subscriptionPaidHandler)

	if err = consumer.StartListening(); err != nil {
		log.Fatal(err)
	}
}

func treatmentCreatedHandler(msg []byte) {
	//TODO	Criar Box com todos os packs e todos os medicamentos

	//p := models.Pack{}
	//err := json.Unmarshal(msg, &p)
	//if err != nil {
	//	fmt.Println("[ERROR] ", err.Error())
	//	return
	//}

	//p.From = time.Now()
	//p.To = time.Now().AddDate(0, 1, 0)
	//p.Status = models.PackStatusWaitingPayment

	//db := postgres.GetDatabase(DatabaseConfig)
	//p.Create(db)

	//b, err := json.Marshal(p)
	//if err != nil {
	//	fmt.Println("[ERROR] ", err.Error())
	//}
	//producer.Publish("pack_created", b)
}

func subscriptionPaidHandler(msg []byte) {
	//TODO mudar status da box para scheduled

	//subs := subscriptionModels.Subscription{}
	//err := json.Unmarshal(msg, &subs)
	//if err != nil {
	//	fmt.Println("[ERROR] ", err.Error())
	//	return
	//}

	//db := postgres.GetDatabase(DatabaseConfig)
	//packs, err := models.GetPacksByOwnerAndStatus(subs.Owner, models.PackStatusWaitingPayment, db)
	//if err != nil {
	//	fmt.Println("[ERROR] ", err.Error())
	//	return
	//}

	//for _, pack := range packs {
	//	pack.Status = models.PackStatusScheduled
	//	pack.Save(db)
	//}
}

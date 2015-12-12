package main

import (
	"fmt"
	"time"

	"github.com/asvins/operations/models"
)

/*
*	Daily CRON for shipping boxes.
 */
func startShipCron() {
	fmt.Println("[INFO] Starting ship CRON")
	go func() {
		shipIt()
		for {
			<-time.After(time.Hour * 24)
			shipIt()
		}
	}()
}

func shipIt() {
	b := models.Box{}
	// TODO write correct query

	toShip, err := b.Retrieve()
	if err != nil {
		fmt.Println("[ERROR] Unable to retrieve boxes to ship: ", err.Error())
		return
	}

	// TODO what is shipping?
}

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/asvins/operations/models"
)

/*
*	Daily CRON for shipping boxes.
 */
func startShipCron() {
	fmt.Println("[INFO] Starting ship CRON")
	go func() {
		allIt()
		for {
			<-time.After(time.Hour * 24)
			allIt()
		}
	}()
}

// VERIFICAÇÂO DE ESTOQUE(consumir)
// /api/inventory/product/:id/consume/:quantity
func consumeFromWarehouse(box *models.Box) error {
	for _, pack := range box.Packs {
		for _, medication := range pack.PackMedications {
			medId := strconv.Itoa(medication.MedicationId)
			url := os.Getenv("") + "/api/inventory/product/" + medId + "/consume/1"

			response, err := http.Get(url)
			if err != nil {
				return err
			}

			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				b, _ := ioutil.ReadAll(response.Body)

				return errors.New(string(b))
			}

		}
	}
	return nil
}

/*
*	1) SHIPAR BOX QUE SE INICIAM EM 7 DIAS COM UMA MARGEM DE +-1 dia - ok
 */
func shipIt() {
	b := models.Box{}
	now := time.Now().Unix()
	nowPlus6Dyas := now + 6*24*60*60
	nowPlus8Days := now + 8*24*60*60

	gteSlice := []string{"start_date|" + strconv.FormatInt(nowPlus6Dyas, 64)}
	lteSlice := []string{"start_date|" + strconv.FormatInt(nowPlus8Days, 64)}
	eqSlice := []string{"status|" + strconv.Itoa(models.BOX_PENDING)}

	b.Base.Query["gte"] = gteSlice
	b.Base.Query["lte"] = lteSlice
	b.Base.Query["eq"] = eqSlice

	toShip, err := b.Retrieve(db)
	if err != nil {
		fmt.Println("[ERROR] Unable to retrieve boxes to ship: ", err.Error())
		return
	}
	for _, curr := range toShip {
		err := consumeFromWarehouse(&curr)
		if err != nil {
			fmt.Println("[ERROR] ", err.Error())
		} else {
			if ok := curr.VerifyOwnerPaymentStatus(); ok == true {
				curr.Status = models.BOX_SHIPED
				if err := curr.Update(db); err != nil {
					fmt.Println("[ERROR] ", err.Error())
					return
				}
			}
		}
	}
}

//	3) As box que começam hoje trocar o status de shiped para DELIVERED
func onIt() {
	b := models.Box{}
	now := time.Now().Unix()
	nowPlus1Day := now + 1*24*60*60
	nowMinus1Day := now - 1*24*60*60

	gteSlice := []string{"start_date|" + strconv.FormatInt(nowPlus1Day, 64)}
	lteSlice := []string{"start_date|" + strconv.FormatInt(nowMinus1Day, 64)}
	eqSlice := []string{"status|" + strconv.Itoa(models.BOX_SHIPED)}

	b.Base.Query["gte"] = gteSlice
	b.Base.Query["lte"] = lteSlice
	b.Base.Query["eq"] = eqSlice

	toOn, err := b.Retrieve(db)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return
	}

	for _, curr := range toOn {
		err := consumeFromWarehouse(&curr)
		if err != nil {
			fmt.Println("[ERROR] ", err.Error())
		} else {
			curr.Status = models.BOX_DELIVERED
			if err := curr.Update(db); err != nil {
				fmt.Println("[ERROR] ", err.Error())
				return
			}
		}
	}

}

//	4) As box que terminar ontem trocar o status para OFF
func offIt() {
	b := models.Box{}
	now := time.Now().Unix()
	nowMinus2Days := now - 2*24*60*60

	gteSlice := []string{"end_date|" + strconv.FormatInt(nowMinus2Days, 64)}
	lteSlice := []string{"end_date|" + strconv.FormatInt(now, 64)}
	eqSlice := []string{"status|" + strconv.Itoa(models.BOX_DELIVERED)}

	b.Base.Query["gte"] = gteSlice
	b.Base.Query["lte"] = lteSlice
	b.Base.Query["eq"] = eqSlice

	toOff, err := b.Retrieve(db)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return
	}

	for _, curr := range toOff {
		err := consumeFromWarehouse(&curr)
		if err != nil {
			fmt.Println("[ERROR] ", err.Error())
		} else {
			curr.Status = models.BOX_FINISHED
			if err := curr.Update(db); err != nil {
				fmt.Println("[ERROR] ", err.Error())
				return
			}
		}
	}

}

// Execute all actions for shipping and updating everything
func allIt() {
	shipIt()
	onIt()
	offIt()
}

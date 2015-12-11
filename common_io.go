package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/asvins/common_io"
	coreModels "github.com/asvins/core/models"
	"github.com/asvins/operations/models"
	"github.com/asvins/utils/config"
)

const (
	FOUR_HOURS   = 4 * 60 * 60
	SIX_HOURS    = 6 * 60 * 60
	EIGHT_HOURS  = 8 * 60 * 60
	TWELVE_HOURS = 12 * 60 * 60
	ONE_DAY      = 24 * 60 * 60
	ONE_MONTH    = 30 * 24 * 60 * 60
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
	t := coreModels.Treatment{}
	err := json.Unmarshal(msg, &t)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return
	}

	// 1) packmads
	var ndays int = ((t.FinishDate - t.StartDate) / ONE_DAY) + 1
	fmt.Println("[DEBUG] ndays: ", ndays)

	packMap := make(map[int][]models.PackMedication)
	for _, currPrescr := range t.Prescriptions {
		var increment int

		switch currPrescr.Frequency {
		case coreModels.PRESCRIPTION_FREQ_4H:
			fmt.Println("[DEBUG] FREQ_4H")
			increment = FOUR_HOURS
			break

		case coreModels.PRESCRIPTION_FREQ_6H:
			fmt.Println("[DEBUG] FREQ_6H")
			increment = SIX_HOURS
			break

		case coreModels.PRESCRIPTION_FREQ_8H:
			fmt.Println("[DEBUG] FREQ_8H")
			increment = EIGHT_HOURS
			break

		case coreModels.PRESCRIPTION_FREQ_12H:
			fmt.Println("[DEBUG] FREQ_12H")
			increment = TWELVE_HOURS
			break

		case coreModels.PRESCRIPTION_FREQ_24H:
			fmt.Println("[DEBUG] FREQ_24H")
			increment = ONE_DAY
			break
		}

		currDate := currPrescr.StartingAt
		for currDate < currPrescr.FinishingAt {
			packMap[currDate] = append(packMap[currDate], models.PackMedication{MedicationId: currPrescr.MedicationId, Quantity: 1})
			currDate += increment
		}

	}

	fmt.Println("[DEBUG] STEP 1 MAP SIZE ", len(packMap))
	fmt.Println("[DEBUG] packMap: ", packMap)

	// 2) packs
	packs := []models.Pack{}
	for date, pmeds := range packMap {
		packs = append(packs, models.Pack{Date: date, TrackingCode: generateTrackingCode(), PackMedications: pmeds})
	}

	// 3) Sort packs by date
	sort.Sort(models.ByDate(packs))
	fmt.Println("[DEBUG] after sort map len", len(packs))
	fmt.Println("[DEBUG] packs: ", packs)

	// 4) box
	currBoxFinalDate := t.StartDate + ONE_MONTH
	currBoxPacks := []models.Pack{}

	for _, currPack := range packs {
		if currPack.Date < currBoxFinalDate {
			currBoxPacks = append(currBoxPacks, currPack)
		} else {
			box := models.Box{
				Status:      models.BOX_SCHEDULED,
				StartDate:   currBoxFinalDate - ONE_MONTH,
				EndDate:     currBoxFinalDate,
				TreatmentId: t.ID,
				PatientId:   t.PatientId,
				Packs:       currBoxPacks,
			}
			err := box.Save(db)
			if err != nil {
				fmt.Println("[ERROR] Could not save box on database: ", err.Error())
				return
			}
			currBoxFinalDate += ONE_MONTH
			currBoxPacks = []models.Pack{}
		}
	}

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

/*
*	Helpers
 */
func generateTrackingCode() string {
	rand.Seed(time.Now().UTC().UnixNano())
	h := sha1.New()
	tc := strconv.Itoa(rand.Intn(10000))
	h.Write([]byte(tc))
	return hex.EncodeToString(h.Sum(nil))
}

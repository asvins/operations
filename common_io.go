package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/asvins/common_io"
	coreModels "github.com/asvins/core/models"
	"github.com/asvins/operations/models"
	subscriptionModels "github.com/asvins/subscription/models"
	"github.com/asvins/utils/config"
	warehouseModels "github.com/asvins/warehouse/models"
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

		warehouseMeds, err := getMedicationsFromWarehouse()
		if err != nil {
			fmt.Println("[ERROR] ", err.Error())
			return
		}

		currDate := currPrescr.StartingAt
		for currDate < currPrescr.FinishingAt {
			toAppend := models.PackMedication{
				MedicationId: currPrescr.MedicationId,
				Quantity:     1,
				Value:        warehouseMeds[currPrescr.MedicationId].CurrentValue,
			}
			packMap[currDate] = append(packMap[currDate], toAppend)
			currDate += increment
		}
	}

	fmt.Println("[DEBUG] STEP 1 MAP SIZE ", len(packMap))
	fmt.Println("[DEBUG] packMap: ", packMap)

	// 2) packs
	packs := []models.Pack{}
	for date, pmeds := range packMap {
		packValue := 0.0
		for _, pmed := range pmeds {
			packValue += pmed.Value
		}
		packs = append(packs, models.Pack{Date: date, TrackingCode: generateTrackingCode(), PackMedications: pmeds, Email: t.Email, Value: packValue})
	}

	// 3) Sort packs by date
	sort.Sort(models.ByDate(packs))
	fmt.Println("[DEBUG] after sort map len", len(packs))
	fmt.Println("[DEBUG] packs: ", packs)

	// 4) box
	currBoxFinalDate := t.StartDate + ONE_MONTH
	if t.FinishDate < ONE_MONTH {
		currBoxFinalDate = t.FinishDate
	}

	currBoxPacks := []models.Pack{}

	for _, currPack := range packs {
		if currPack.Date < currBoxFinalDate {
			currBoxPacks = append(currBoxPacks, currPack)
		} else {
			//createBox()
			currBoxValue := 0.0
			for _, currBoxPacksPack := range currBoxPacks {
				currBoxValue += currBoxPacksPack.Value
			}
			createBox(currBoxPacks, currBoxFinalDate, t)
			currBoxFinalDate += ONE_MONTH
			currBoxPacks = []models.Pack{}
		}
	}
	if len(currBoxPacks) != 0 {
		createBox(currBoxPacks, currBoxFinalDate, t)
	}
}

func createBox(currBoxPacks []models.Pack, currBoxFinalDate int, t coreModels.Treatment) {
	currBoxValue := 0.0
	for _, currBoxPacksPack := range currBoxPacks {
		currBoxValue += currBoxPacksPack.Value
	}

	fmt.Println("[INFO] Will save box with ", len(currBoxPacks), " Packs!")
	box := models.Box{
		Status:      models.BOX_PENDING,
		StartDate:   currBoxFinalDate - ONE_MONTH,
		EndDate:     currBoxFinalDate,
		TreatmentId: t.ID,
		PatientId:   t.PatientId,
		Packs:       currBoxPacks,
		Value:       currBoxValue,
	}
	err := box.Save(db)
	if err != nil {
		fmt.Println("[ERROR] Could not save box on database: ", err.Error())
		return
	}
	sendBoxCreated(box)
}

func subscriptionPaidHandler(msg []byte) {
	subs := subscriptionModels.Subscription{}

	if err := json.Unmarshal(msg, &subs); err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return
	}

	// 1) GET no core para pegar os tratamentos do paciente com status = INATIVO
	baseUrl := "http://" + os.Getenv("DEPLOY_CORE_1_PORT_8080_TCP_ADDR") + ":" + os.Getenv("DEPLOY_CORE_1_PORT_8080_TCP_PORT")
	resp, err := http.Get(baseUrl + "/api/treatments?eq=patient_id|" + subs.Owner + "&eq=status|" + strconv.Itoa(coreModels.TREATMENT_STATUS_INACTIVE))
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return
	}

	ts := []coreModels.Treatment{}

	if err := json.Unmarshal(body, &ts); err != nil {
		fmt.Println("[ERROR] ", err.Error())
	}
	// Até aqui foi só para montar a struct do tratamento

	for _, currTreatment := range ts {
		box := models.Box{}
		box.TreatmentId = currTreatment.ID
		box.Status = models.BOX_PENDING
		boxes, err := box.RetrieveOrdered(db)
		if err != nil {
			fmt.Println("[ERROR] ", err.Error())
			return
		}

		if len(boxes) > 0 {
			boxes[0].Status = models.BOX_SCHEDULED
			boxes[0].Update(db)
		} else {
			fmt.Println("[INFO] No boxes to schedule")
		}
	}

	b, err := json.Marshal(ts)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return
	}

	producer.Publish("activate_treatments", b)
}

/*
*	Senders
 */
func sendBoxCreated(box models.Box) {
	topic, _ := common_io.BuildTopicFromCommonEvent(common_io.EVENT_CREATED, "box")
	/*
	 * json Marshal
	 */
	b, err := json.Marshal(&box)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return
	}

	producer.Publish(topic, b)
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

func getMedicationsFromWarehouse() (map[int]warehouseModels.Product, error) {
	baseURL := "http://" + os.Getenv("DEPLOY_WAREHOUSE_1_PORT_8080_TCP_ADDR") + ":" + os.Getenv("DEPLOY_WAREHOUSE_1_PORT_8080_TCP_PORT")
	response, err := http.Get(baseURL + "/api/inventory/product")
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	products := []warehouseModels.Product{}
	if err := json.Unmarshal(body, &products); err != nil {
		return nil, err
	}

	m := make(map[int]warehouseModels.Product)
	for _, p := range products {
		m[p.ID] = p
	}

	return m, nil
}

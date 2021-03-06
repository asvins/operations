package models

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
)

const (
	BOX_PENDING = iota
	BOX_SCHEDULED
	BOX_SHIPED
	BOX_DELIVERED
	BOX_FINISHED
)

type Box struct {
	Base
	ID          int     `json:"id"`
	Status      int     `json:"status"`
	StartDate   int64   `json:"start_date"`
	EndDate     int64   `json:"end_date"`
	TreatmentId int     `json:"treatment_id"`
	PatientId   int     `json:"patient_id"`
	Value       float64 `json:"value"`
	Packs       []Pack  `json:"packs"`
}

func (b *Box) Save(db *gorm.DB) error {
	return db.Create(b).Error
}

func (b *Box) Update(db *gorm.DB) error {
	return db.Save(b).Error
}

func (b *Box) Delete(db *gorm.DB) error {
	return db.Delete(b).Error
}

func (b *Box) Retrieve(db *gorm.DB) ([]Box, error) {
	var bs []Box
	err := db.Where(b).Find(&bs, b.Base.BuildQuery()).Error

	return bs, err
}

func (b *Box) RetrieveOrdered(db *gorm.DB) ([]Box, error) {
	var bs []Box
	err := db.Order("start_date asc").Where(b).Find(&bs, b.Base.BuildQuery()).Error

	return bs, err
}

func (b *Box) VerifyOwnerPaymentStatus() bool {
	resp, err := http.Get(os.Getenv("") + "/api/subscriber/" + strconv.Itoa(b.PatientId) + "/paymentstatus")
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true
	} else {
		return false
	}
}

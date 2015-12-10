package models

import "github.com/jinzhu/gorm"

type PackMedication struct {
	Base
	ID           int `json:"id"`
	MedicationId int `json:"medication_id"`
	PackId       int `json:"pack_id"`
	Quantity     int `json:"quantity"`
}

func (p *PackMedication) Save(db *gorm.DB) error {
	return db.Create(p).Error
}

func (p *PackMedication) Update(db *gorm.DB) error {
	return db.Save(p).Error
}

func (p *PackMedication) Delete(db *gorm.DB) error {
	return db.Delete(p).Error
}

func (p *PackMedication) Retrieve(db *gorm.DB) ([]PackMedication, error) {
	var ps []PackMedication
	err := db.Where(p).Find(&ps, p.Base.BuildQuery()).Error

	return ps, err
}

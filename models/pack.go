package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Pack struct {
	Base
	ID              int `json:"pack_id" gorm:"column:id"`
	BoxId           int
	Date            time.Time        `json:"from" gorm:"column:from_date"`
	TrackingCode    string           `json:"tracking_code"`
	PackMedications []PackMedication `json:"pack_medications"`
}

func (p *Pack) Save(db *gorm.DB) error {
	return db.Create(p).Error
}

func (p *Pack) Update(db *gorm.DB) error {
	return db.Save(p).Error
}

func (p *Pack) Delete(db *gorm.DB) error {
	return db.Delete(p).Error
}

func (p *Pack) Retrieve(db *gorm.DB) ([]Pack, error) {
	var ps []Pack
	err := db.Where(p).Find(&ps, p.Base.BuildQuery()).Error

	return ps, err
}

package models

import "github.com/jinzhu/gorm"

type Pack struct {
	Base
	ID              int              `json:"pack_id"`
	BoxId           int              `json:"box_id"`
	Date            int              `json:"date"`
	TrackingCode    string           `json:"tracking_code"`
	PackMedications []PackMedication `json:"pack_medications"`
}

/*
*	Sort interface implementation
 */
type ByDate []Pack

func (ps ByDate) Len() int {
	return len(ps)
}

func (ps ByDate) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps ByDate) Less(i, j int) bool {
	return ps[i].Date < ps[j].Date
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

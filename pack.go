package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	PackStatusDelivered = iota
	PackStatusShipped
	PackStatusOnProduction
	PackStatusScheduled
	PackStatusWaitingPayment
)

type Pack struct {
	ID           int       `json:"pack_id" gorm:"column:id"`
	Owner        string    `json:"owner"`
	Supervisor   string    `json:"supervisor"`
	From         time.Time `json:"from" gorm:"column:from_date"`
	To           time.Time `json:"to" gorm:"column:to_date"`
	TrackingCode string    `json:"tracking_code"`
	Status       int       `json:"Status"`
	PackType     string    `json:"pack_type"`
	PackHash     string    `json:"hash"`
}

func statusToString(status int) string {
	switch status {
	case PackStatusDelivered:
		return "delivered"
	case PackStatusShipped:
		return "shipped"
	case PackStatusOnProduction:
		return "production"
	case PackStatusScheduled:
		return "scheduled"
	}
	return ""
}

func stringToStatus(status string) int {
	switch status {
	case "delivered":
		return PackStatusDelivered
	case "shipped":
		return PackStatusShipped
	case "production":
		return PackStatusOnProduction
	case "scheduled":
		return PackStatusScheduled
	}
	return 100000
}

func (p *Pack) Create(db *gorm.DB) error {
	return db.Create(p).Error
}

func (p *Pack) Save(db *gorm.DB) error {
	return db.Save(p).Error
}

func GetPacksByOwnerAndStatus(owner string, status int, db *gorm.DB) ([]Pack, error) {
	var packs []Pack
	err := db.Where("owner = ? and status = ?", owner, status).Find(&packs).Error
	return packs, err
}

func GetPacksByOwner(owner string, ps *[]Pack, db *gorm.DB) error {
	return db.Where("owner = ?", owner).Find(ps).Error
}

func GetPacksByStatus(status int, ps *[]Pack, db *gorm.DB) error {
	return db.Where("status = ?", status).Find(ps).Error
}

func GetPacksByStatusString(status string, ps *[]Pack, db *gorm.DB) error {
	return GetPacksByStatus(stringToStatus(status), ps, db)
}

func GetActivePacks(owner string, ps *[]Pack, db *gorm.DB) error {
	twoWeeksAgo := time.Now().AddDate(0, 0, -14)
	return db.Where("owner = ? AND to_date > ?", owner, twoWeeksAgo).Find(ps).Error
}

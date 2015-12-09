package models

import "time"

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

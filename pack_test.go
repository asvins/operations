package main

import (
	"testing"
	"time"

	"github.com/asvins/common_db/postgres"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetPacks(t *testing.T) {
	p1 := &Pack{Owner: "john.doe@example.com", Supervisor: "admin@asvins.com", From: time.Now().AddDate(0, 1, 0), To: time.Now().AddDate(0, 0, -15), TrackingCode: "A29CDE", Status: PackStatusDelivered, PackType: "medication", PackHash: "ccf3b89a"}
	p2 := &Pack{Owner: "john.doe@example.com", Supervisor: "admin@asvins.com", From: time.Now(), To: time.Now().AddDate(0, 0, 14), TrackingCode: "A29CDE", Status: PackStatusShipped, PackType: "medication", PackHash: "ccf3b89a"}
	db := postgres.GetDatabase(DatabaseConfig)
	p1.Create(db)
	p2.Create(db)
	Convey("When getting packs", t, func() {
		Convey("We can get them by owner", func() {
			var p []Pack
			GetPacksByOwner("john.doe@example.com", &p, db)
			So(len(p), ShouldEqual, 2)
		})
		Convey("We can get them by status", func() {
			var p []Pack
			GetPacksByStatus(PackStatusDelivered, &p, db)
			So(len(p), ShouldEqual, 1)
			GetPacksByStatusString("shipped", &p, db)
			So(len(p), ShouldEqual, 1)
		})
		Convey("We can get active packs", func() {
			var p []Pack
			GetActivePacks("john.doe@example.com", &p, db)
			So(len(p), ShouldEqual, 1)
		})
	})
	db.Exec("TRUNCATE TABLE packs")
}

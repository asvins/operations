package main

import (
	"net/http"

	"github.com/asvins/operations/models"
	"github.com/asvins/router/errors"
)

func retrieveBoxes(w http.ResponseWriter, r *http.Request) errors.Http {
	b := models.Box{}
	if err := BuildStructFromQueryString(&b, r.URL.Query()); err != nil {
		return errors.BadRequest(err.Error())
	}

	b.Base.Query = r.URL.Query()

	boxes, err := b.Retrieve(db)
	if err != nil {
		return errors.BadRequest(err.Error())
	}

	if len(boxes) == 0 {
		return errors.NotFound("record not found")
	}
	rend.JSON(w, http.StatusOK, boxes)

	return nil
}

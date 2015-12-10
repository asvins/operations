package main

import (
	"net/http"

	"github.com/asvins/operations/models"
	"github.com/asvins/router/errors"
)

func retrievePacks(w http.ResponseWriter, r *http.Request) errors.Http {
	p := models.Pack{}
	if err := BuildStructFromQueryString(&p, r.URL.Query()); err != nil {
		return errors.BadRequest(err.Error())
	}

	p.Base.Query = r.URL.Query()

	packs, err := p.Retrieve(db)
	if err != nil {
		return errors.BadRequest(err.Error())
	}

	if len(packs) == 0 {
		return errors.NotFound("record not found")
	}
	rend.JSON(w, http.StatusOK, packs)

	return nil
}

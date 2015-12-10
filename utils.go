package main

import (
	"io"
	"net/url"

	"github.com/asvins/warehouse/decoder"
)

func BuildStructFromQueryString(dst interface{}, queryString url.Values) error {
	decoder := decoder.NewDecoder()
	return decoder.DecodeURLValues(dst, queryString)
}

func BuildStructFromReqBody(dst interface{}, body io.ReadCloser) error {
	decoder := decoder.NewDecoder()
	return decoder.DecodeReqBody(dst, body)
}

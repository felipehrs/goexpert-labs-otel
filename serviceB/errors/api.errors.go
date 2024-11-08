package errors

import "errors"

var (
	InvalidZipCode          = errors.New("invalid zipcode")
	NotFoundZipCode         = errors.New("can not find zipcode")
	UnableToRetrieveWeather = errors.New("unable to retrieve weather")
	UnableToRetrieveZipCode = errors.New("unable to retrieve zipcode")
)

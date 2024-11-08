package usecase

import (
	apiErrors "github.com/felipehrs/goexpert-labs-otel-serciceA/errors"
)

type WeatherUsecase interface {
	GetWeatherByCep(zipcode string) (WeatherResponse, error)
}

type weatherUsecase struct {
}

type ZipCodeInput struct {
	ZipCode string `json:"cep"`
}

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func NewWeatherUsecase() WeatherUsecase {
	return &weatherUsecase{}
}

func (w *weatherUsecase) GetWeatherByCep(zipCode string) (WeatherResponse, error) {
	return WeatherResponse{}, apiErrors.InvalidZipCode
}

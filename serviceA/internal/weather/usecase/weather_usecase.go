package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	apiErrors "github.com/felipehrs/goexpert-labs-otel-serciceA/errors"
	"github.com/felipehrs/goexpert-labs-otel-serciceA/pkg"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

	if !pkg.IsValidZipCode(zipCode) {
		return WeatherResponse{}, apiErrors.InvalidZipCode
	}

	//TODO: Call serviceB to get weather by zipCode

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	url := fmt.Sprintf("http://localhost:8080/weather/%s", zipCode)

	resp, err := client.Get(url)

	if err != nil {
		return WeatherResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return WeatherResponse{}, apiErrors.NotFoundZipCode
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherResponse{}, err
	}

	var response WeatherResponse
	err = json.Unmarshal(body, &response)

	if err != nil {
		return WeatherResponse{}, err
	}

	return response, nil

}

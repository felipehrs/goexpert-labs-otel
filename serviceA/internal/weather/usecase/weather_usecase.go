package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	apiErrors "github.com/felipehrs/goexpert-labs-otel-serciceA/errors"
	pkg "github.com/felipehrs/goexpert-labs-otel-serciceA/pkg"
)

type WeatherUsecase interface {
	GetWeatherByCep(zipcode string) (WeatherResponse, error)
}

type weatherUsecase struct {
}

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ViaCepResponse struct {
	CEP        string `json:"cep"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	Estado     string `json:"estado"`
	Erro       string `json:"erro"`
}

type OpenWeatherResponse struct {
	Current Current `json:"current"`
}

type Current struct {
	TempC float64 `json:"temp_c"`
}

func NewWeatherUsecase() WeatherUsecase {
	return &weatherUsecase{}
}

func (w *weatherUsecase) GetWeatherByCep(zipCode string) (WeatherResponse, error) {
	if zipCode == "" {
		return WeatherResponse{}, apiErrors.InvalidZipCode
	}

	if !pkg.IsValidZipCode(zipCode) {
		return WeatherResponse{}, apiErrors.InvalidZipCode
	}

	viaCepResponse, err := w.doViaCepRequest(zipCode)
	if errors.Is(err, apiErrors.NotFoundZipCode) {
		return WeatherResponse{}, apiErrors.NotFoundZipCode
	}

	if err != nil {
		return WeatherResponse{}, apiErrors.UnableToRetrieveZipCode
	}

	weatherApiResponse, err := w.doWeatherRequest(viaCepResponse.Localidade)
	if err != nil {
		return WeatherResponse{}, apiErrors.UnableToRetrieveWeather
	}

	return WeatherResponse{
		TempC: weatherApiResponse.Current.TempC,
		TempF: CelsiusToFahrenheit(weatherApiResponse.Current.TempC),
		TempK: CelsiusToKelvin(weatherApiResponse.Current.TempC),
	}, nil

}

func (w *weatherUsecase) doViaCepRequest(zipCode string) (ViaCepResponse, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", zipCode)

	resp, err := http.Get(url)
	if err != nil {
		return ViaCepResponse{}, err
	}
	defer resp.Body.Close()

	var viaCepResponse ViaCepResponse
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &viaCepResponse); err != nil {
		return ViaCepResponse{}, err
	}

	if viaCepResponse.Erro != "" {
		return ViaCepResponse{}, apiErrors.NotFoundZipCode
	}

	return viaCepResponse, nil
}

func (w *weatherUsecase) doWeatherRequest(city string) (OpenWeatherResponse, error) {
	apiKey := "5903bf6c2eb54a7f863170134240311"
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, city)

	resp, err := http.Get(url)

	if err != nil {
		return OpenWeatherResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OpenWeatherResponse{}, err
	}

	var weatherApiResponse OpenWeatherResponse
	err = json.Unmarshal(body, &weatherApiResponse)

	if err != nil {
		return OpenWeatherResponse{}, err
	}

	return weatherApiResponse, nil
}

func CelsiusToFahrenheit(celsius float64) float64 {
	return celsius*9/5 + 32
}

func CelsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}

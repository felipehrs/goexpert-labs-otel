package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"io"
	"net/http"
	"net/url"
	"unicode"

	apiErrors "github.com/felipehrs/goexpert-labs-otel-serciceB/errors"
	pkg "github.com/felipehrs/goexpert-labs-otel-serciceB/pkg"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

type WeatherUsecase interface {
	GetWeatherByCep(c *gin.Context, zipcode string) (WeatherResponse, error)
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

func (w *weatherUsecase) GetWeatherByCep(c *gin.Context, zipCode string) (WeatherResponse, error) {
	if zipCode == "" {
		return WeatherResponse{}, apiErrors.InvalidZipCode
	}

	if !pkg.IsValidZipCode(zipCode) {
		return WeatherResponse{}, apiErrors.InvalidZipCode
	}

	ctx := c.Request.Context()
	tracer := otel.Tracer("weatherUsecase.GetWeatherByCep")
	ctx, span := tracer.Start(ctx, "CEP_REQUEST")
	defer span.End()

	viaCepResponse, err := w.doViaCepRequest(zipCode)
	if errors.Is(err, apiErrors.NotFoundZipCode) {
		return WeatherResponse{}, apiErrors.NotFoundZipCode
	}

	if err != nil {
		return WeatherResponse{}, apiErrors.UnableToRetrieveZipCode
	}

	tracer = otel.Tracer("weatherUsecase.GetWeatherByCep")
	_, span = tracer.Start(ctx, "WEATHER_REQUEST")
	defer span.End()

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

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, city)

	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, url.QueryEscape(result))

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

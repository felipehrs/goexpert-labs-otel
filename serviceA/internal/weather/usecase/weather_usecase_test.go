package usecase_test

import (
	"testing"

	apiErrors "github.com/felipehrs/goexpert-labs-otel-serciceA/errors"
	. "github.com/felipehrs/goexpert-labs-otel-serciceA/internal/weather/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var InvalidZipCode = "123456789"
var ValidZipCode = "89201440"
var NotFoundZipCode = "00000000"
var EmptyCep = ""

type WeatherUsecaseTestSuite struct {
	suite.Suite
	weatherUsecase WeatherUsecase
}

func (s *WeatherUsecaseTestSuite) SetupSuite() {
	s.weatherUsecase = NewWeatherUsecase()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(WeatherUsecaseTestSuite))
}

func (s *WeatherUsecaseTestSuite) TestGetWeatherByCep_EmptyZipCode() {
	_, err := s.weatherUsecase.GetWeatherByCep(EmptyCep)
	s.Equal(apiErrors.InvalidZipCode, err)
}

func (s *WeatherUsecaseTestSuite) TestGetWeatherByCep_InvalidZipCode() {
	_, err := s.weatherUsecase.GetWeatherByCep(InvalidZipCode)
	s.Equal(apiErrors.InvalidZipCode, err)
}

func (s *WeatherUsecaseTestSuite) TestGetWeatherByCep_NotFound() {
	_, err := s.weatherUsecase.GetWeatherByCep(NotFoundZipCode)
	s.Equal(apiErrors.NotFoundZipCode, err)
}

func (s *WeatherUsecaseTestSuite) TestGetWeatherByCep_ValidZipCode() {
	weather, err := s.weatherUsecase.GetWeatherByCep(ValidZipCode)
	s.NoError(err)
	s.NotEmpty(weather)
}

func (s *WeatherUsecaseTestSuite) TestShould_ReturnWeatherByCep() {
	weather, err := s.weatherUsecase.GetWeatherByCep(ValidZipCode)
	s.NoError(err)
	s.NotEmpty(weather)
	s.NotZero(weather.TempC)
	s.NotZero(weather.TempF)
	s.NotZero(weather.TempK)
}

func TestCelsiusToFahrenheit(t *testing.T) {
	assert.Equal(t, 77.0, CelsiusToFahrenheit(25))
	assert.Equal(t, 32.0, CelsiusToFahrenheit(0))
	assert.Equal(t, 212.0, CelsiusToFahrenheit(100))
}

func TestCelsiusToKelvin(t *testing.T) {
	assert.Equal(t, 298.0, CelsiusToKelvin(25))
	assert.Equal(t, 273.0, CelsiusToKelvin(0))
	assert.Equal(t, 373.0, CelsiusToKelvin(100))
}

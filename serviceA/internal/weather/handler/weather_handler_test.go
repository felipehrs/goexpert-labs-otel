package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	apiError "github.com/felipehrs/goexpert-labs-otel-serciceA/errors"
	"github.com/felipehrs/goexpert-labs-otel-serciceA/internal/weather/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	. "github.com/felipehrs/goexpert-labs-otel-serciceA/internal/weather/handler"
)

type MockWeatherUsecase struct {
	mock.Mock
}

func (m *MockWeatherUsecase) GetWeatherByCep(zipcode string) (usecase.WeatherResponse, error) {
	args := m.Called(zipcode)
	if args.Get(0) != nil {
		return args.Get(0).(usecase.WeatherResponse), args.Error(1)
	}
	return usecase.WeatherResponse{}, args.Error(1)
}

func TestWeatherHandler_Handle_InvalidZipCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(MockWeatherUsecase)
	mockUsecase.On("GetWeatherByCep", "1234").Return(usecase.WeatherResponse{}, apiError.InvalidZipCode)

	handler := NewWeatherHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "zipcode", Value: "1234"})

	handler.Handle(c)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	exectBody := `{"error":"invalid zipcode"}`
	assert.JSONEq(t, exectBody, w.Body.String())

	mockUsecase.AssertExpectations(t)
}

func TestWeatherHandler_Handle_ZipCodeNotFount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(MockWeatherUsecase)
	mockUsecase.On("GetWeatherByCep", "00000000").Return(usecase.WeatherResponse{}, apiError.NotFoundZipCode)

	handler := NewWeatherHandler(mockUsecase)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "zipcode", Value: "00000000"})

	handler.Handle(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

	exectBody := `{"error":"can not find zipcode"}`

	assert.JSONEq(t, exectBody, w.Body.String())

	mockUsecase.AssertExpectations(t)
}

func TestWeatherHandler_Handle_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(MockWeatherUsecase)
	mockUsecase.On("GetWeatherByCep", "89201215").Return(usecase.WeatherResponse{
		TempC: 25.0,
		TempF: 77.0,
		TempK: 298.15,
	}, nil)

	handler := NewWeatherHandler(mockUsecase)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Params = append(c.Params, gin.Param{Key: "zipcode", Value: "89201215"})

	handler.Handle(c)

	assert.Equal(t, http.StatusOK, w.Code)

	exectBody := `{"temp_C":25,"temp_F":77,"temp_K":298.15}`

	assert.JSONEq(t, exectBody, w.Body.String())

	mockUsecase.AssertExpectations(t)
}

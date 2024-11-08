package handler

import (
	"errors"
	"net/http"

	"github.com/felipehrs/goexpert-labs-otel-serciceB/internal/weather/usecase"
	"github.com/gin-gonic/gin"

	apiErrors "github.com/felipehrs/goexpert-labs-otel-serciceB/errors"
)

type WeatherHandler interface {
	Handle(ctx *gin.Context)
}

type weatherHandler struct {
	usecase usecase.WeatherUsecase
}

func NewWeatherHandler(usecase usecase.WeatherUsecase) WeatherHandler {
	return &weatherHandler{usecase: usecase}
}

func (w *weatherHandler) Handle(ctx *gin.Context) {
	zipCode := ctx.Param("zipcode")

	weather, err := w.usecase.GetWeatherByCep(zipCode)

	if errors.Is(err, apiErrors.InvalidZipCode) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": apiErrors.InvalidZipCode.Error()})
		return
	}

	if errors.Is(err, apiErrors.NotFoundZipCode) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": apiErrors.NotFoundZipCode.Error()})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, weather)
}

package main

import (
	"github.com/felipehrs/goexpert-labs-otel-serciceB/internal/weather/handler"
	"github.com/felipehrs/goexpert-labs-otel-serciceB/internal/weather/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	usecase := usecase.NewWeatherUsecase()
	handler := handler.NewWeatherHandler(usecase)

	r.GET("/weather/:zipcode", handler.Handle)
	r.Run(":8080")
}

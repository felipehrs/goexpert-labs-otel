package main

import (
	"github.com/felipehrs/goexpert-labs-otel-serciceA/internal/weather/handler"
	"github.com/felipehrs/goexpert-labs-otel-serciceA/internal/weather/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	usecase := usecase.NewWeatherUsecase()
	handler := handler.NewWeatherHandler(usecase)

	r.POST("/weather", handler.Handle)
	r.Run(":8080")
}

package main

import (
	"github.com/RuslanBoo/metrics/internal/handler"
	"github.com/RuslanBoo/metrics/internal/repository/memory"
	"github.com/RuslanBoo/metrics/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	metricStorage := memory.New()
	metricService := service.New(metricStorage)
	metricHandler := handler.New(metricService)

	metricRouter := gin.Default()
	metricHandler.RegisterRoutes(metricRouter)

	if err := metricRouter.Run(":8080"); err != nil {
		panic(err)
	}
}

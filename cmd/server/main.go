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
	metricConf := service.ParseServerFlags()

	metricRouter := gin.Default()
	metricHandler.RegisterRoutes(metricRouter)

	if err := metricRouter.Run(metricConf.Addr); err != nil {
		panic(err)
	}
}

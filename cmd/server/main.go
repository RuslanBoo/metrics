package main

import (
	"net/http"

	"github.com/RuslanBoo/metrics/internal/handler"
	"github.com/RuslanBoo/metrics/internal/repository/memory"
	"github.com/RuslanBoo/metrics/internal/service"
)

func main() {
	metricStorage := memory.New()
	metricService := service.New(metricStorage)
	metricHandler := handler.New(metricService)

	metricMux := http.NewServeMux()
	metricMux.HandleFunc(`/update/`, metricHandler.Update)

	err := http.ListenAndServe(`:8080`, metricMux)
	if err != nil {
		panic(err)
	}
}

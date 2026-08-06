package main

import (
	"log"

	"featherweight/internal/config"
	deliveryhttp "featherweight/internal/delivery/http"
	"featherweight/internal/repository/memory"
)

func main() {

	//load configuration
	cfg := config.Load()

	log.Printf("starting server on %s", cfg.ServerAddress)

	jobRepository := memory.NewJobRepository()
	log.Printf("Job repository initialized: %T", jobRepository)

	//initialize gin router

	router := deliveryhttp.NewRouter(cfg)

	//start server

	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

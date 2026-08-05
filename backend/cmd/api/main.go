package main

import (
	"featherweight/internal/config"
	deliveryhttp "featherweight/internal/delivery/http"
	"log"
)

func main() {

	//load configuration
	cfg := config.Load()

	log.Printf("starting server on %s", cfg.ServerAddress)

	//initialize gin router

	router := deliveryhttp.NewRouter(cfg)

	//start server

	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

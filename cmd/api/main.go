package main

import (
	"fmt"
	"log"
	"net/http"

	"chainhub-api/internal/config"
	"chainhub-api/internal/db"
	apihttp "chainhub-api/internal/http"
	"chainhub-api/internal/http/handlers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	defer dbConn.Close()

	handler := handlers.New(dbConn, cfg)
	router := apihttp.NewRouter(handler)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("listening on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

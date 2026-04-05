package main

import (
	"log"
	"net/http"
	"os"
	"purr-case/internal/httpapi"
	"purr-case/internal/httpapi/global"
	"purr-case/internal/httpapi/items"
	"purr-case/internal/httpapi/payments"
	"purr-case/internal/httpapi/users"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	merchant_id := os.Getenv("merchant_id")

	gh := global.InitHandler()
	uh := users.InitHandler()
	ih := items.InitHandler(merchant_id)
	ph := payments.InitHandler()
	router := httpapi.NewRouter(gh, uh, ih, ph)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Starting server on port %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

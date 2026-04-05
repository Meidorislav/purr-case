package main

import (
	"log"
	"net/http"
	"os"
	"purr-case/internal/httpapi"
	"purr-case/internal/httpapi/global"
	"purr-case/internal/httpapi/payments"
	"purr-case/internal/httpapi/users"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	gh := global.InitHandler()
	ph := payments.InitHandler()
	uh := users.InitHandler()
	router := httpapi.NewRouter(gh, ph, uh)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Starting server on port %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

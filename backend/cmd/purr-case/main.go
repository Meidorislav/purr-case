package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"purr-case/internal/db"
	"purr-case/internal/httpapi"
	"purr-case/internal/httpapi/global"
	"purr-case/internal/httpapi/inventory"
	"purr-case/internal/httpapi/items"
	"purr-case/internal/httpapi/payments"
	"purr-case/internal/httpapi/users"
	inventory_service "purr-case/internal/service/inventory"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	merchant_id := os.Getenv("merchant_id")
	xsollaSandbox := strings.EqualFold(os.Getenv("XSOLLA_SANDBOX"), "true")
	xsollaProjectID, _ := strconv.Atoi(os.Getenv("XSOLLA_PROJECT_ID"))
	xsollaAPIKey := os.Getenv("XSOLLA_API_KEY")
	xsollaReturnURL := os.Getenv("XSOLLA_RETURN_URL")
	xsollaSandbox := strings.EqualFold(os.Getenv("XSOLLA_SANDBOX"), "true")
	xsollaProjectID, _ := strconv.Atoi(os.Getenv("XSOLLA_PROJECT_ID"))
	xsollaAPIKey := os.Getenv("XSOLLA_API_KEY")
	xsollaReturnURL := os.Getenv("XSOLLA_RETURN_URL")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	database, err := db.InitDatabase(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer database.Pool.Close()

	log.Println("Successfully connected to database!")

	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("migrations connected")

	is := inventory_service.InitService(database)

	gh := global.InitHandler()
	uh := users.InitHandler()
	ih := items.InitHandler(merchant_id)
	invh := inventory.InitHandler(is)
	ph := payments.InitHandler(payments.Config{
		MerchantID: merchant_id,
		ProjectID:  xsollaProjectID,
		APIKey:     xsollaAPIKey,
		ReturnURL:  xsollaReturnURL,
		Sandbox:    xsollaSandbox,
	})
	router := httpapi.NewRouter(gh, uh, ih, ph, invh)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	}()

	log.Printf("Starting server on port %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

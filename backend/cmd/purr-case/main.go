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
	"strconv"
	"strings"
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

	gh := global.InitHandler()
	uh := users.InitHandler()
	ih := items.InitHandler(merchant_id)
	ph := payments.InitHandler(payments.Config{
		MerchantID: merchant_id,
		ProjectID:  xsollaProjectID,
		APIKey:     xsollaAPIKey,
		ReturnURL:  xsollaReturnURL,
		Sandbox:    xsollaSandbox,
	})
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

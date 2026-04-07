package payments

import (
	"bytes"
	"context"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	dto "purr-case/internal/dto/payments"
	"purr-case/internal/httpapi/respond"
	inventory_service "purr-case/internal/service/inventory"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// Config holds Xsolla credentials and behavioural flags for the payments handler.
// All values are read from environment variables at startup (see main.go).
type Config struct {
	MerchantID       string // Xsolla merchant ID used for Basic auth and token API URL.
	ProjectID        int    // Xsolla project ID embedded in every token request.
	APIKey           string // Xsolla API key (Basic auth password).
	WebhookSecretKey string // Shared secret used to verify incoming webhook signatures.
	ReturnURL        string // URL Xsolla redirects the user to after payment completes.
	Sandbox          bool   // When true, uses Xsolla sandbox endpoints and sets mode="sandbox".
}

// xsollaUserField wraps a single scalar value as required by the Xsolla token API schema.
type xsollaUserField struct {
	Value string `json:"value"`
}

// xsollaUserCountry carries the user's billing country.
// AllowModify=false prevents the user from changing it on the Pay Station page.
type xsollaUserCountry struct {
	Value       string `json:"value"`
	AllowModify bool   `json:"allow_modify"`
}

type xsollaTokenUser struct {
	ID      xsollaUserField   `json:"id"`
	Country xsollaUserCountry `json:"country,omitempty"`
}

type xsollaTokenSettings struct {
	ExternalID string        `json:"external_id,omitempty"`
	Language   string        `json:"language,omitempty"`
	ReturnURL  string        `json:"return_url,omitempty"`
	UI         xsollaTokenUI `json:"ui,omitempty"`
}

type xsollaTokenUI struct {
	Theme               string          `json:"theme,omitempty"`
	IsCartOpenByDefault bool            `json:"is_cart_open_by_default,omitempty"`
	Desktop             xsollaUIDesktop `json:"desktop,omitempty"`
}

type xsollaUIDesktop struct {
	Header xsollaUIHeader `json:"header,omitempty"`
}

type xsollaUIHeader struct {
	IsVisible       bool   `json:"is_visible,omitempty"`
	VisibleLogo     bool   `json:"visible_logo,omitempty"`
	VisibleName     bool   `json:"visible_name,omitempty"`
	VisiblePurchase bool   `json:"visible_purchase,omitempty"`
	Type            string `json:"type,omitempty"`
	CloseButton     bool   `json:"close_button,omitempty"`
	CloseButtonIcon string `json:"close_button_icon,omitempty"`
}

type xsollaPurchaseItem struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

// xsollaTokenRequest is the body sent to the Xsolla token API.
// The returned token is used to build the Pay Station checkout URL.
type xsollaTokenRequest struct {
	Sandbox  bool            `json:"sandbox,omitempty"`
	User     xsollaTokenUser `json:"user"`
	Purchase struct {
		Items []xsollaPurchaseItem `json:"items"`
	} `json:"purchase"`
	Settings xsollaTokenSettings `json:"settings"`
}

// xsollaWebhookPayload covers the common fields shared across all Xsolla webhook events.
// Xsolla sends different subsets of fields depending on notification_type,
// so most nested fields are optional and may be zero-valued.
type xsollaWebhookPayload struct {
	NotificationType string `json:"notification_type"`
	User             struct {
		ID      any    `json:"id"` // can be string or number depending on Xsolla project settings
		Email   string `json:"email"`
		Country string `json:"country"`
	} `json:"user"`
	Payment struct {
		Status string `json:"status"`
	} `json:"payment"`
	Transaction struct {
		ID         any    `json:"id"`
		ExternalID string `json:"external_id"` // our internal order ID if passed during token creation
	} `json:"transaction"`
	Order struct {
		ID         any    `json:"id"`
		ExternalID string `json:"external_id"`
		Status     string `json:"status"`
	} `json:"order"`
}

type Handler struct {
	cfg       Config
	inventory *inventory_service.Service
}

func InitHandler(cfg Config, inventory *inventory_service.Service) *Handler {
	return &Handler{cfg: cfg, inventory: inventory}
}

// CreateCheckout creates an Xsolla payment token for the given cart and returns
// a Pay Station URL that the client should redirect the user to.
// The user ID is taken from the JWT set by the Auth middleware.
func (h *Handler) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCheckoutRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get user ID from JWT token set by the Auth middleware
	userID, ok := r.Context().Value("userId").(string)
	if !ok || userID == "" {
		respond.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if len(req.Items) == 0 {
		respond.WriteError(w, http.StatusBadRequest, "items are required")
		return
	}

	seenSKUs := make(map[string]struct{}, len(req.Items))

	// Validate only the fields that should come from the storefront.
	// Pricing, currency and the final order amount must be resolved by Xsolla.
	for i, item := range req.Items {
		if item.SKU == "" {
			respond.WriteError(w, http.StatusBadRequest, "items["+strconv.Itoa(i)+"].sku is required")
			return
		}
		if item.Quantity <= 0 {
			respond.WriteError(w, http.StatusBadRequest, "items["+strconv.Itoa(i)+"].quantity must be a positive integer")
			return
		}

		if _, exists := seenSKUs[item.SKU]; exists {
			respond.WriteError(w, http.StatusBadRequest, "duplicate sku values are not supported in one checkout request")
			return
		}
		seenSKUs[item.SKU] = struct{}{}
	}

	orderID := "order-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	xsollaReq := h.buildXsollaTokenRequest(userID, orderID, req)

	xsollaResp, err := h.createXsollaToken(xsollaReq)
	if err != nil {
		respond.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	if err := h.storeCheckoutOrder(r.Context(), userID, orderID, req.Items); err != nil {
		respond.WriteError(w, http.StatusInternalServerError, "failed to store checkout order")
		return
	}

	resp := dto.CreateCheckoutResponse{
		OrderID:     orderID,
		Status:      "new",
		ItemsCount:  len(req.Items),
		CheckoutURL: h.buildCheckoutURL(xsollaResp.Token),
	}

	respond.WriteJSON(w, http.StatusOK, resp)
}

// HandleWebhook receives and validates webhook notifications from Xsolla.
// Xsolla signs every request with SHA1(body + secret); the signature is verified
// before the payload is processed. Responds with 200 on success or a structured
// Xsolla-compatible error body on failure.
//
// NOTE: post-payment business logic (e.g. granting items to the user's inventory)
// is not implemented here — that is the responsibility of the inventory team.
func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		respond.WriteError(w, http.StatusBadRequest, "failed to read webhook payload")
		return
	}

	if len(rawBody) == 0 {
		respond.WriteError(w, http.StatusBadRequest, "empty webhook payload")
		return
	}

	if h.cfg.WebhookSecretKey == "" {
		respond.WriteError(w, http.StatusInternalServerError, "xsolla webhook secret is not configured")
		return
	}

	if !verifyXsollaWebhookSignature(rawBody, h.cfg.WebhookSecretKey, r.Header.Get("Authorization")) {
		writeXsollaWebhookError(w, http.StatusBadRequest, "INVALID_SIGNATURE", "Invalid signature")
		return
	}

	var payload xsollaWebhookPayload
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		respond.WriteError(w, http.StatusBadRequest, "invalid webhook payload")
		return
	}

	if payload.NotificationType == "" {
		respond.WriteError(w, http.StatusBadRequest, "webhook event type is required")
		return
	}

	userID := stringifyWebhookValue(payload.User.ID)
	if payload.NotificationType == "user_validation" && userID == "" {
		writeXsollaWebhookError(w, http.StatusBadRequest, "INVALID_USER", "Invalid user")
		return
	}

	eventStatus := h.resolveWebhookStatus(payload)
	externalID := firstNonEmpty(payload.Order.ExternalID, payload.Transaction.ExternalID)
	orderID := firstNonEmpty(externalID, stringifyWebhookValue(payload.Order.ID))
	transactionID := firstNonEmpty(stringifyWebhookValue(payload.Transaction.ID), payload.Transaction.ExternalID)

	fulfilled := false
	if payload.NotificationType == "order_paid" {
		if externalID == "" {
			respond.WriteError(w, http.StatusBadRequest, "webhook external id is required")
			return
		}

		fulfilled, err = h.fulfillPaidOrder(r.Context(), externalID, transactionID)
		if err != nil {
			respond.WriteError(w, http.StatusInternalServerError, "failed to fulfill paid order")
			return
		}
	}

	respond.WriteJSON(w, http.StatusOK, map[string]any{
		"received":         true,
		"fulfilled":        fulfilled,
		"notificationType": payload.NotificationType,
		"status":           eventStatus,
		"userId":           userID,
		"orderId":          orderID,
		"transactionId":    transactionID,
	})
}

// buildCheckoutURL constructs the Pay Station URL from a payment token.
// Uses the sandbox domain when Config.Sandbox is true.
func (h *Handler) buildCheckoutURL(token string) string {
	if h.cfg.Sandbox {
		return "https://sandbox-secure.xsolla.com/paystation4/?token=" + token
	}

	return "https://secure.xsolla.com/paystation4/?token=" + token
}

func (h *Handler) buildXsollaTokenRequest(userID string, orderID string, req dto.CreateCheckoutRequest) xsollaTokenRequest {
	settings := xsollaTokenSettings{
		ExternalID: orderID,
		Language:   "en",
		ReturnURL:  h.cfg.ReturnURL,
		UI: xsollaTokenUI{
			Theme:               "63295aab2e47fab76f7708e3",
			IsCartOpenByDefault: true,
			Desktop: xsollaUIDesktop{
				Header: xsollaUIHeader{
					IsVisible:       true,
					VisibleLogo:     true,
					VisibleName:     true,
					VisiblePurchase: true,
					Type:            "normal",
					CloseButton:     true,
					CloseButtonIcon: "cross",
				},
			},
		},
	}

	items := make([]xsollaPurchaseItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, xsollaPurchaseItem{
			SKU:      item.SKU,
			Quantity: item.Quantity,
		})
	}

	payload := xsollaTokenRequest{
		Sandbox: h.cfg.Sandbox,
		User: xsollaTokenUser{
			ID: xsollaUserField{
				Value: userID,
			},
			// Temporary fallback required by Xsolla to choose order currency.
			// Later this should come from the real user profile or X-User-Ip.
			Country: xsollaUserCountry{
				Value:       "US",
				AllowModify: false,
			},
		},
		Settings: settings,
	}
	payload.Purchase.Items = items

	return payload
}

type xsollaTokenResponse struct {
	Token string `json:"token"`
}

// createXsollaToken calls the Xsolla token API and returns a short-lived payment token.
// Authentication uses HTTP Basic with project ID as username and API key as password.
func (h *Handler) createXsollaToken(payload xsollaTokenRequest) (xsollaTokenResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return xsollaTokenResponse{}, fmt.Errorf("marshal xsolla payload: %w", err)
	}

	url := fmt.Sprintf(
		"https://store.xsolla.com/api/v3/project/%d/admin/payment/token",
		h.cfg.ProjectID,
	)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return xsollaTokenResponse{}, fmt.Errorf("create xsolla request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(h.cfg.ProjectID) + ":" + h.cfg.APIKey))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return xsollaTokenResponse{}, fmt.Errorf("send xsolla request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return xsollaTokenResponse{}, fmt.Errorf("xsolla returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var xsollaResp xsollaTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&xsollaResp); err != nil {
		return xsollaTokenResponse{}, fmt.Errorf("decode xsolla response: %w", err)
	}

	if xsollaResp.Token == "" {
		return xsollaTokenResponse{}, fmt.Errorf("xsolla token is empty")
	}

	return xsollaResp, nil
}

func (h *Handler) storeCheckoutOrder(ctx context.Context, userID string, orderID string, items []dto.CheckoutItem) error {
	if h.inventory == nil || h.inventory.Database == nil {
		return fmt.Errorf("inventory service is not configured")
	}

	// Store the checkout before the user pays, so the webhook can later map
	// Xsolla's external_id back to the exact SKUs and quantities we sold.
	tx, err := h.inventory.Database.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin store checkout transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO payment_orders (external_id, user_id, status)
		 VALUES ($1, $2, 'new')
		 ON CONFLICT (external_id) DO NOTHING`,
		orderID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("insert payment order: %w", err)
	}

	for _, item := range items {
		_, err = tx.Exec(ctx,
			`INSERT INTO payment_order_items (external_id, sku, quantity)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (external_id, sku)
			 DO UPDATE SET quantity = EXCLUDED.quantity`,
			orderID,
			item.SKU,
			item.Quantity,
		)
		if err != nil {
			return fmt.Errorf("insert payment order item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit store checkout transaction: %w", err)
	}

	return nil
}

func (h *Handler) fulfillPaidOrder(ctx context.Context, externalID string, transactionID string) (bool, error) {
	if h.inventory == nil || h.inventory.Database == nil {
		return false, fmt.Errorf("inventory service is not configured")
	}

	// Prefer the Xsolla transaction ID for idempotency. Fall back to externalID
	// for manual tests or webhook shapes that do not include transaction.id.
	eventKey := firstNonEmpty(transactionID, externalID)

	tx, err := h.inventory.Database.Pool.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("begin fulfill order transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var userID string
	if err := tx.QueryRow(ctx,
		`SELECT user_id FROM payment_orders WHERE external_id = $1 FOR UPDATE`,
		externalID,
	).Scan(&userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("payment order %s not found", externalID)
		}
		return false, fmt.Errorf("select payment order: %w", err)
	}

	// Xsolla can retry webhooks. This insert is the guard that prevents us
	// from granting the same paid order twice.
	result, err := tx.Exec(ctx,
		`INSERT INTO processed_payment_events (event_key, external_id)
		 VALUES ($1, $2)
		 ON CONFLICT (event_key) DO NOTHING`,
		eventKey,
		externalID,
	)
	if err != nil {
		return false, fmt.Errorf("insert processed payment event: %w", err)
	}
	if result.RowsAffected() == 0 {
		if err := tx.Commit(ctx); err != nil {
			return false, fmt.Errorf("commit duplicate payment event transaction: %w", err)
		}
		return false, nil
	}

	// The webhook only tells us which order was paid. The list of items comes
	// from the checkout snapshot saved when we created the Xsolla token.
	rows, err := tx.Query(ctx,
		`SELECT sku, quantity FROM payment_order_items WHERE external_id = $1`,
		externalID,
	)
	if err != nil {
		return false, fmt.Errorf("query payment order items: %w", err)
	}
	defer rows.Close()

	var items []inventory_service.GrantItem
	for rows.Next() {
		var item inventory_service.GrantItem
		if err := rows.Scan(&item.SKU, &item.Quantity); err != nil {
			return false, fmt.Errorf("scan payment order item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("payment order item rows error: %w", err)
	}
	if len(items) == 0 {
		return false, fmt.Errorf("payment order %s has no items", externalID)
	}

	// Actual inventory mutation. It runs inside the same DB transaction as the
	// idempotency insert, so a failure cannot leave the order half-processed.
	if err := h.inventory.GrantItemsInTx(ctx, tx, userID, items); err != nil {
		return false, fmt.Errorf("grant paid order items: %w", err)
	}

	// Mark the local order as fulfilled after the inventory upsert succeeds.
	_, err = tx.Exec(ctx,
		`UPDATE payment_orders
		 SET status = 'paid',
		     transaction_id = NULLIF($2, ''),
		     processed_at = COALESCE(processed_at, now())
		 WHERE external_id = $1`,
		externalID,
		transactionID,
	)
	if err != nil {
		return false, fmt.Errorf("update payment order status: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("commit fulfill order transaction: %w", err)
	}

	return true, nil
}

// resolveWebhookStatus maps an Xsolla notification type to a human-readable status string.
// Falls back to the nested payment/order status when available.
func (h *Handler) resolveWebhookStatus(payload xsollaWebhookPayload) string {
	switch payload.NotificationType {
	case "user_validation":
		return "validated"
	case "payment":
		return firstNonEmpty(payload.Payment.Status, "paid")
	case "refund":
		return firstNonEmpty(payload.Payment.Status, "refunded")
	case "order_paid":
		return firstNonEmpty(payload.Order.Status, payload.Payment.Status, "paid")
	case "order_canceled":
		return firstNonEmpty(payload.Order.Status, payload.Payment.Status, "canceled")
	default:
		return firstNonEmpty(payload.Order.Status, payload.Payment.Status, "received")
	}
}

// verifyXsollaWebhookSignature checks that the Authorization header matches
// the expected Xsolla signature: SHA1(rawBody + secretKey) in hex, prefixed with "Signature ".
// Uses constant-time comparison to prevent timing attacks.
func verifyXsollaWebhookSignature(rawBody []byte, secretKey string, authorizationHeader string) bool {
	if secretKey == "" {
		return false
	}

	signature := strings.TrimSpace(authorizationHeader)
	if signature == "" {
		return false
	}

	signature = strings.TrimPrefix(signature, "Signature ")
	signature = strings.TrimSpace(signature)
	if signature == "" {
		return false
	}

	sum := sha1.Sum(append(rawBody, []byte(secretKey)...))
	computed := fmt.Sprintf("%x", sum)
	received := strings.ToLower(signature)

	return subtle.ConstantTimeCompare([]byte(computed), []byte(received)) == 1
}

// stringifyWebhookValue converts Xsolla webhook ID fields to string.
// Xsolla may send user/order/transaction IDs as either a JSON string or number,
// so we normalise them here to avoid type assertion errors downstream.
func stringifyWebhookValue(v any) string {
	switch typed := v.(type) {
	case nil:
		return ""
	case string:
		return typed
	case float64:
		return strconv.FormatInt(int64(typed), 10)
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	case json.Number:
		return typed.String()
	default:
		return fmt.Sprintf("%v", typed)
	}
}

// firstNonEmpty returns the first non-blank string from the provided values.
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}

// writeXsollaWebhookError writes an error response in the format expected by Xsolla.
// Xsolla requires a specific JSON shape to recognise and display webhook errors correctly.
func writeXsollaWebhookError(w http.ResponseWriter, status int, code string, message string) {
	respond.WriteJSON(w, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

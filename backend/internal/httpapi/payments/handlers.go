package payments

import (
	"bytes"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	dto "purr-case/internal/dto/payments"
	"purr-case/internal/httpapi/respond"
	"strconv"
	"strings"
	"time"
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
	ExternalID string `json:"external_id,omitempty"`
	Language   string `json:"language,omitempty"`
	ReturnURL  string `json:"return_url,omitempty"`
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
	cfg Config
}

func InitHandler(cfg Config) *Handler {
	return &Handler{cfg: cfg}
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
	orderID := firstNonEmpty(payload.Order.ExternalID, stringifyWebhookValue(payload.Order.ID))
	transactionID := firstNonEmpty(payload.Transaction.ExternalID, stringifyWebhookValue(payload.Transaction.ID))

	respond.WriteJSON(w, http.StatusOK, map[string]any{
		"received":         true,
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

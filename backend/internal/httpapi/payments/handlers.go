package payments

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	dto "purr-case/internal/dto/payments"
	"purr-case/internal/httpapi/respond"
	"strconv"
	"time"
)

type Config struct {
	MerchantID string
	ProjectID  int
	APIKey     string
	ReturnURL  string
	Sandbox    bool
}

type xsollaUserField struct {
	Value string `json:"value"`
}

type xsollaUserCountry struct {
	Value       string `json:"value"`
	AllowModify bool   `json:"allow_modify"`
}

type xsollaTokenUser struct {
	ID      xsollaUserField   `json:"id"`
	Country xsollaUserCountry `json:"country,omitempty"`
}

type xsollaTokenSettings struct {
	ProjectID  int    `json:"project_id"`
	ExternalID string `json:"external_id,omitempty"`
	Language   string `json:"language,omitempty"`
	ReturnURL  string `json:"return_url,omitempty"`
	Mode       string `json:"mode,omitempty"` // "sandbox" for test payments
}

type xsollaTokenRequest struct {
	User     xsollaTokenUser     `json:"user"`
	Settings xsollaTokenSettings `json:"settings"`
}

type Handler struct {
	cfg Config
}

func InitHandler(cfg Config) *Handler {
	return &Handler{cfg: cfg}
}

// Entry point for creating a payment and obtaining the payment link
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

// Reserved entry point for handling webhooks from the payment provider
func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var payload map[string]any

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respond.WriteError(w, http.StatusBadRequest, "invalid webhook payload")
		return
	}

	if payload["event"] == nil && payload["notification_type"] == nil && payload["type"] == nil {
		respond.WriteError(w, http.StatusBadRequest, "webhook event type is required")
		return
	}

	respond.WriteJSON(w, http.StatusOK, map[string]any{
		"message":  "webhook received",
		"received": true,
	})
}

func (h *Handler) buildCheckoutURL(token string) string {
	if h.cfg.Sandbox {
		return "https://sandbox-secure.xsolla.com/paystation4/?token=" + token
	}

	return "https://secure.xsolla.com/paystation4/?token=" + token
}

func (h *Handler) buildXsollaTokenRequest(userID string, orderID string, req dto.CreateCheckoutRequest) xsollaTokenRequest {
	settings := xsollaTokenSettings{
		ProjectID:  h.cfg.ProjectID,
		ExternalID: orderID,
		Language:   "en",
		ReturnURL:  h.cfg.ReturnURL,
	}
	if h.cfg.Sandbox {
		settings.Mode = "sandbox"
	}

	return xsollaTokenRequest{
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
}

type xsollaTokenResponse struct {
	Token string `json:"token"`
}

func (h *Handler) createXsollaToken(payload xsollaTokenRequest) (xsollaTokenResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return xsollaTokenResponse{}, fmt.Errorf("marshal xsolla payload: %w", err)
	}

	url := fmt.Sprintf(
		"https://api.xsolla.com/merchant/v2/merchants/%s/token",
		h.cfg.MerchantID,
	)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return xsollaTokenResponse{}, fmt.Errorf("create xsolla request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(h.cfg.MerchantID + ":" + h.cfg.APIKey))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return xsollaTokenResponse{}, fmt.Errorf("send xsolla request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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

package cases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	dto "purr-case/internal/dto/cases"

	"purr-case/internal/db"
	inventory_service "purr-case/internal/service/inventory"
)

var ErrNotInInventory = errors.New("case not in inventory")
var ErrCaseNotFound = errors.New("case not found")
var ErrInvalidDropTable = errors.New("case has no valid drop table")

type Service struct {
	db            *db.Database
	invSvc        *inventory_service.Service
	xsollaBaseURL string
}

func InitService(database *db.Database, invSvc *inventory_service.Service, projectID string) *Service {
	return &Service{
		db:            database,
		invSvc:        invSvc,
		xsollaBaseURL: fmt.Sprintf("https://store.xsolla.com/api/v2/project/%s", projectID),
	}
}

func (s *Service) fetchDropTable(ctx context.Context, caseSKU, token string) ([]dto.DropEntry, error) {
	url := fmt.Sprintf("%s/items/bundle/sku/%s?additional_fields[]=custom_attributes", s.xsollaBaseURL, caseSKU)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch bundle: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrCaseNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("xsolla returned status %d", resp.StatusCode)
	}

	var bundle struct {
		CustomAttributes json.RawMessage `json:"custom_attributes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&bundle); err != nil {
		return nil, fmt.Errorf("decode bundle response: %w", err)
	}

	var attrs dto.CaseCustomAttributes
	if err := json.Unmarshal(bundle.CustomAttributes, &attrs); err != nil {
		return nil, fmt.Errorf("decode custom_attributes: %w", err)
	}

	if len(attrs.DropTable) == 0 {
		return nil, ErrInvalidDropTable
	}

	return attrs.DropTable, nil
}

func rollItem(table []dto.DropEntry) string {
	total := 0
	for _, e := range table {
		total += e.Weight
	}
	r := rand.Intn(total)
	for _, e := range table {
		r -= e.Weight
		if r < 0 {
			return e.SKU
		}
	}
	return table[len(table)-1].SKU
}

func (s *Service) OpenCase(ctx context.Context, userID, caseSKU, token string) (string, error) {
	dropTable, err := s.fetchDropTable(ctx, caseSKU, token)
	if err != nil {
		return "", err
	}

	wonSKU := rollItem(dropTable)

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx,
		`UPDATE inventory SET quantity = quantity - 1 WHERE user_id = $1 AND sku = $2 AND quantity > 0`,
		userID, caseSKU,
	)
	if err != nil {
		return "", fmt.Errorf("decrement case: %w", err)
	}
	if result.RowsAffected() == 0 {
		return "", ErrNotInInventory
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM inventory WHERE user_id = $1 AND sku = $2 AND quantity <= 0`,
		userID, caseSKU,
	)
	if err != nil {
		return "", fmt.Errorf("cleanup zero quantity: %w", err)
	}

	if err := s.invSvc.GrantItemsInTx(ctx, tx, userID, []inventory_service.GrantItem{{SKU: wonSKU, Quantity: 1}}); err != nil {
		return "", fmt.Errorf("grant won item: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit transaction: %w", err)
	}

	return wonSKU, nil
}

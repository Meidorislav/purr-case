package catalog_service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	dto "purr-case/internal/dto/items"
)

type Service struct {
	BaseURL string
	Client  *http.Client
}

func InitService(projectID string) *Service {
	return &Service{
		BaseURL: fmt.Sprintf("https://store.xsolla.com/api/v2/project/%s", projectID),
		Client:  http.DefaultClient,
	}
}

func (s *Service) FetchItems(ctx context.Context, token string, itemType string, rawQuery string) (dto.CatalogResponse, error) {
	requestURL := s.BaseURL + "/items" + itemType
	if rawQuery != "" {
		requestURL += "?" + rawQuery
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return dto.CatalogResponse{}, fmt.Errorf("create catalog request: %w", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return dto.CatalogResponse{}, fmt.Errorf("fetch catalog items: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return dto.CatalogResponse{}, fmt.Errorf("fetch catalog items: unexpected status %d", resp.StatusCode)
	}

	var result dto.CatalogResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return dto.CatalogResponse{}, fmt.Errorf("decode catalog response: %w", err)
	}

	return result, nil
}

func (s *Service) GetCatalogItems(ctx context.Context, token string) ([]dto.Item, error) {
	itemTypes := []string{"", "/virtual_items"}
	itemsBySKU := make(map[string]dto.Item)

	for _, itemType := range itemTypes {
		result, err := s.FetchItems(ctx, token, itemType, "")
		if err != nil {
			return nil, err
		}

		for _, item := range result.Items {
			if item.SKU == "" {
				continue
			}
			itemsBySKU[item.SKU] = item
		}
	}

	items := make([]dto.Item, 0, len(itemsBySKU))
	for _, item := range itemsBySKU {
		items = append(items, item)
	}

	return items, nil
}

func (s *Service) FetchItemBySKU(ctx context.Context, token string, sku string, rawQuery string) (dto.Item, error) {
	requestURL := s.BaseURL + "/items/sku/" + url.PathEscape(sku)
	if rawQuery != "" {
		requestURL += "?" + rawQuery
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return dto.Item{}, fmt.Errorf("create catalog item request: %w", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return dto.Item{}, fmt.Errorf("fetch catalog item: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return dto.Item{}, fmt.Errorf("fetch catalog item: unexpected status %d", resp.StatusCode)
	}

	var result dto.Item
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return dto.Item{}, fmt.Errorf("decode catalog item response: %w", err)
	}

	return result, nil
}

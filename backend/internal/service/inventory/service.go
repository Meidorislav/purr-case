package inventory_service

import (
	"context"
	"fmt"

	"purr-case/internal/db"
	dto "purr-case/internal/dto/inventory"
)

type Service struct {
	Database *db.Database
}

func InitService(db *db.Database) *Service {
	return &Service{
		Database: db,
	}
}

func (s *Service) GetUserInventory(ctx context.Context, userID string) ([]dto.InventoryItem, error) {
	rows, err := s.Database.Pool.Query(ctx,
		`SELECT id, user_id, sku, quantity FROM inventory WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("query inventory: %w", err)
	}
	defer rows.Close()

	var items []dto.InventoryItem
	for rows.Next() {
		var item dto.InventoryItem
		if err := rows.Scan(&item.ID, &item.UserID, &item.SKU, &item.Quantity); err != nil {
			return nil, fmt.Errorf("scan inventory row: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return items, nil
}

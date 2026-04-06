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

func (s *Service) UpdateUserInventoryItem(ctx context.Context, userID string, sku string, quantity int) (bool, error) {
	result, err := s.Database.Pool.Exec(ctx,
		`UPDATE inventory SET quantity = $1 WHERE user_id = $2 AND sku = $3`,
		quantity, userID, sku,
	)
	if err != nil {
		return false, fmt.Errorf("update inventory item: %w", err)
	}

	return result.RowsAffected() > 0, nil
}

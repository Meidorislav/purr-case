package inventory_service

import (
	"context"
	"fmt"

	"purr-case/internal/db"
	dto "purr-case/internal/dto/inventory"

	"github.com/jackc/pgx/v5"
)

type Service struct {
	Database *db.Database
}

type GrantItem struct {
	SKU      string
	Quantity int
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

func (s *Service) GrantItems(ctx context.Context, userID string, items []GrantItem) error {
	tx, err := s.Database.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin grant items transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.GrantItemsInTx(ctx, tx, userID, items); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit grant items transaction: %w", err)
	}

	return nil
}

func (s *Service) GrantItemsInTx(ctx context.Context, tx pgx.Tx, userID string, items []GrantItem) error {
	for _, item := range items {
		if item.SKU == "" || item.Quantity <= 0 {
			return fmt.Errorf("invalid grant item")
		}

		_, err := tx.Exec(ctx,
			`INSERT INTO inventory (user_id, sku, quantity)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (user_id, sku)
			 DO UPDATE SET quantity = inventory.quantity + EXCLUDED.quantity`,
			userID,
			item.SKU,
			item.Quantity,
		)
		if err != nil {
			return fmt.Errorf("grant inventory item: %w", err)
		}
	}

	return nil
}

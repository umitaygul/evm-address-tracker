package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/umitaygul/evm-address-tracker/internal/models"
)

type AddressRepository struct {
	db *pgxpool.Pool
}

func NewAddressRepository(db *pgxpool.Pool) *AddressRepository {
	return &AddressRepository{db: db}
}

func (r *AddressRepository) Create(ctx context.Context, userID string, chainID int64, address string) (*models.WatchedAddress, error) {
	query := `
		INSERT INTO watched_addresses (user_id, chain_id, address)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, chain_id, address, created_at`

	var a models.WatchedAddress
	err := r.db.QueryRow(ctx, query, userID, chainID, address).Scan(
		&a.ID, &a.UserID, &a.ChainID, &a.Address, &a.CreatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrDuplicate
		}
		return nil, err
	}
	return &a, nil
}

func (r *AddressRepository) ListByUser(ctx context.Context, userID string) ([]models.WatchedAddress, error) {
	query := `
		SELECT id, user_id, chain_id, address, created_at
		FROM watched_addresses
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []models.WatchedAddress
	for rows.Next() {
		var a models.WatchedAddress
		if err := rows.Scan(&a.ID, &a.UserID, &a.ChainID, &a.Address, &a.CreatedAt); err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	if addresses == nil {
		addresses = []models.WatchedAddress{}
	}
	return addresses, rows.Err()
}

func (r *AddressRepository) Delete(ctx context.Context, id, userID string) error {
	query := `DELETE FROM watched_addresses WHERE id = $1 AND user_id = $2`

	tag, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *AddressRepository) GetByID(ctx context.Context, id, userID string) (*models.WatchedAddress, error) {
	query := `
		SELECT id, user_id, chain_id, address, created_at
		FROM watched_addresses
		WHERE id = $1 AND user_id = $2`

	var a models.WatchedAddress
	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&a.ID, &a.UserID, &a.ChainID, &a.Address, &a.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &a, nil
}

package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/umitaygul/evm-address-tracker/internal/models"
)

type WebhookRepository struct {
	db *pgxpool.Pool
}

func NewWebhookRepository(db *pgxpool.Pool) *WebhookRepository {
	return &WebhookRepository{db: db}
}

func (r *WebhookRepository) Create(ctx context.Context, userID, url, secret string) (*models.Webhook, error) {
	query := `
		INSERT INTO webhooks (user_id, url, secret)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, url, secret, created_at`
	var w models.Webhook
	err := r.db.QueryRow(ctx, query, userID, url, secret).Scan(
		&w.ID, &w.UserID, &w.URL, &w.Secret, &w.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WebhookRepository) ListByUser(ctx context.Context, userID string) ([]models.Webhook, error) {
	query := `SELECT id, user_id, url, secret, created_at FROM webhooks WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []models.Webhook
	for rows.Next() {
		var w models.Webhook
		if err := rows.Scan(&w.ID, &w.UserID, &w.URL, &w.Secret, &w.CreatedAt); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, w)
	}
	return webhooks, nil
}

func (r *WebhookRepository) Delete(ctx context.Context, id int64, userID string) error {
	query := `DELETE FROM webhooks WHERE id = $1 AND user_id = $2`
	ct, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *WebhookRepository) GetByID(ctx context.Context, id int64) (*models.Webhook, error) {
	query := `SELECT id, user_id, url, secret, created_at FROM webhooks WHERE id = $1`
	var w models.Webhook
	err := r.db.QueryRow(ctx, query, id).Scan(
		&w.ID, &w.UserID, &w.URL, &w.Secret, &w.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &w, nil
}

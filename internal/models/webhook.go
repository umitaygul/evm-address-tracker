package models

import "time"

type Webhook struct {
	ID        int64     `json:"id"`
	UserID    string    `json:"user_id"`
	URL       string    `json:"url"`
	Secret    string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

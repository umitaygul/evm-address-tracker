package models

import "time"

type WatchedAddress struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ChainID   int64     `json:"chain_id"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

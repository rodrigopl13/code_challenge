package entities

import "time"

type Message struct {
	ID        int       `json:"id"`
	User      User      `json:"user"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

package lsd

import "time"

type Markdown string

type Notification struct {
	ID          string    `json:"id"`
	DeliveredAt time.Time `json:"delivered_at"`
	Title       string    `json:"title"`
	Body        Markdown  `json:"body"`
	Options     []string  `json:"options"`
}

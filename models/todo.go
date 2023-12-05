package models

import (
	"time"
)

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      bool      `json:"status" db:"status"`
	CreatedDate time.Time `json:"created_date" db:"created_date"`
}

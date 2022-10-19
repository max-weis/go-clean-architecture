package entity

import "time"

type Product struct {
	ID          string
	Title       string
	Description string
	Price       uint64
	CreatedAt   time.Time
	ModifiedAt  time.Time
}

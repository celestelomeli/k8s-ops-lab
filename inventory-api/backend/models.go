// models.go defines the data structures used 
package main

import "time"

// Product represents one item in the inventory
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
git s
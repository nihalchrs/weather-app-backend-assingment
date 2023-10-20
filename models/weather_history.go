package models

import "time"

type Weather struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id" binding:"required"`
	City      string    `json:"city"`
	Temp      float64   `json:"temp"`
	Pressure  int       `json:"pressure"`
	Humidity  int       `json:"humidity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

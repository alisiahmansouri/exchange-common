package model

import "github.com/google/uuid"

// model/order.go

type OrderCreateRequest struct {
	UserID        uuid.UUID `json:"user_id" binding:"required"`
	PairID        uuid.UUID `json:"pair_id" binding:"required"`
	Side          string    `json:"side" binding:"required,oneof=buy sell"`
	OrderType     string    `json:"order_type" binding:"required,oneof=limit market"`
	Amount        float64   `json:"amount" binding:"required,gt=0"`
	Price         float64   `json:"price" binding:"required,gte=0"`
	ClientOrderID *string   `json:"client_order_id,omitempty"`
	TimeInForce   *string   `json:"time_in_force,omitempty"`
	Meta          *string   `json:"meta,omitempty"`
}

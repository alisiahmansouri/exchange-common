package model

import (
	"github.com/alisiahmansouri/exchange-common/entity"
	"github.com/google/uuid"
	"time"
)

// ------------------- OrderCreateRequest -------------------

type OrderCreateRequest struct {
	UserID        uuid.UUID `json:"user_id" binding:"required"` // توسط هندلر ست میشه، نیازی نیست کاربر ست کنه
	PairID        uuid.UUID `json:"pair_id" binding:"required"`
	Side          string    `json:"side" binding:"required,oneof=buy sell"`
	OrderType     string    `json:"order_type" binding:"required,oneof=limit market"`
	Amount        float64   `json:"amount" binding:"required,gt=0"`
	Price         float64   `json:"price" binding:"required,gte=0"`
	ClientOrderID *string   `json:"client_order_id,omitempty"`
	TimeInForce   *string   `json:"time_in_force,omitempty"`
	Meta          *string   `json:"meta,omitempty"`
}

// ------------------- OrderResponse (برای خروجی API) -------------------

type OrderResponse struct {
	ID            uuid.UUID  `json:"id"`
	PairID        uuid.UUID  `json:"pair_id"`
	OrderType     string     `json:"order_type"`
	Side          string     `json:"side"`
	Amount        float64    `json:"amount"`
	FilledAmount  float64    `json:"filled_amount"`
	Price         float64    `json:"price"`
	Status        string     `json:"status"`
	TimeInForce   string     `json:"time_in_force"`
	ClientOrderID *string    `json:"client_order_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	ExecutedAt    *time.Time `json:"executed_at,omitempty"`
}

func ToOrderResponse(order *entity.Order) *OrderResponse {
	if order == nil {
		return nil
	}
	return &OrderResponse{
		ID:            order.ID,
		PairID:        order.PairID,
		OrderType:     string(order.OrderType),
		Side:          string(order.Side),
		Amount:        order.Amount,
		FilledAmount:  order.FilledAmount,
		Price:         order.Price,
		Status:        string(order.Status),
		TimeInForce:   string(order.TimeInForce),
		ClientOrderID: order.ClientOrderID,
		CreatedAt:     order.CreatedAt,
		ExecutedAt:    order.ExecutedAt,
	}
}

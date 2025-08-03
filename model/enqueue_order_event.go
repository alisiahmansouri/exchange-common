package model

import "github.com/alisiahmansouri/exchange-common/entity"

type EnqueueOrderEvent struct {
	Order entity.Order `json:"order"`
}

type CancelOrderEvent struct {
	OrderID string `json:"order_id"`
}

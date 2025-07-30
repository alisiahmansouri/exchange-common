package model

type OrderBook struct {
	PairID          string
	Bids            []OrderBookItem
	Asks            []OrderBookItem
	PricePrecision  uint  // دقت قیمت (مثلاً 2)
	AmountPrecision uint  // دقت مقدار (مثلاً 8)
	LastUpdate      int64 // Optional: timestamp آخرین بروزرسانی
}

type OrderBookItem struct {
	Price    float64  `json:"price"`
	Amount   float64  `json:"amount"`
	OrderIDs []string `json:"order_ids,omitempty"`
}

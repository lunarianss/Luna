package biz_entity

type PriceInfo struct {
	UnitPrice   float64 `json:"unit_price"`
	Unit        float64 `json:"unit"`
	TotalAmount float64 `json:"total_amount"`
	Currency    string  `json:"currency"`
}

func NewFreePriceInfo() *PriceInfo {
	return &PriceInfo{
		UnitPrice:   0,
		Unit:        0,
		TotalAmount: 0,
		Currency:    "USD",
	}
}

type PriceType string

const (
	INPUT  PriceType = "input"
	OUTPUT PriceType = "output"
)

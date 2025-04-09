package model

type TradePlateResult struct {
	Direction    string            `json:"direction"`
	MaxAmount    float64           `json:"maxAmount"`
	MinAmount    float64           `json:"minAmount"`
	HighestPrice float64           `json:"highestPrice"`
	LowestPrice  float64           `json:"lowestPrice"`
	Symbol       string            `json:"symbol"`
	Items        []*TradePlateItem `json:"items"`
}
type TradePlateItem struct {
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

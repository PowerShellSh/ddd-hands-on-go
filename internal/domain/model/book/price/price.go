package price

import "fmt"

// Currency は通貨を表す型です。
type Currency string

const (
	// JPY は日本円を表します。
	JPY Currency = "JPY"
)

// Price は価格を表す値オブジェクトです。
type Price struct {
	amount   float64
	currency Currency
}

// NewPrice は新しいPriceを生成します。
// 金額が負の値の場合、または通貨がJPY以外の場合はエラーを返します。
func NewPrice(amount float64, currency Currency) (*Price, error) {
	if amount < 0 {
		return nil, fmt.Errorf("Priceの生成に失敗しました: 金額は0以上である必要があります")
	}
	if currency != JPY {
		return nil, fmt.Errorf("Priceの生成に失敗しました: 通貨はJPYのみサポートされています")
	}
	return &Price{amount: amount, currency: currency}, nil
}

// Amount は金額を返します。
func (p *Price) Amount() float64 {
	return p.amount
}

// Currency は通貨を返します。
func (p *Price) Currency() Currency {
	return p.currency
}

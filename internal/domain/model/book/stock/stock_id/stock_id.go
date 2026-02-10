package stock_id

import "fmt"

// StockId は在庫IDを表す値オブジェクトです。
type StockId struct {
	value string
}

// NewStockId は新しいStockIdを生成します。
// 値が空の場合はエラーを返します。
func NewStockId(value string) (*StockId, error) {
	if value == "" {
		return nil, fmt.Errorf("StockIdの生成に失敗しました: IDは必須です")
	}
	return &StockId{value: value}, nil
}

// Value はStockIdの値を返します。
func (id *StockId) Value() string {
	return id.value
}

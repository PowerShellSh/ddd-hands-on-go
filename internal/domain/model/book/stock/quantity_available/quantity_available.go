package quantity_available

import "fmt"

// QuantityAvailable は在庫数を表す値オブジェクトです。
type QuantityAvailable struct {
	value int
}

// NewQuantityAvailable は新しいQuantityAvailableを生成します。
// 数値が負の場合はエラーを返します。
func NewQuantityAvailable(value int) (*QuantityAvailable, error) {
	if value < 0 {
		return nil, fmt.Errorf("QuantityAvailableの生成に失敗しました: 数値は0以上である必要があります")
	}
	return &QuantityAvailable{value: value}, nil
}

// Value は在庫数を返します。
func (q *QuantityAvailable) Value() int {
	return q.value
}

// Increment は在庫数を増加させた新しいQuantityAvailableを返します。
func (q *QuantityAvailable) Increment(amount int) (*QuantityAvailable, error) {
	if amount < 0 {
		return nil, fmt.Errorf("負の値で増加させることはできません")
	}
	return NewQuantityAvailable(q.value + amount)
}

// Decrement は在庫数を減少させた新しいQuantityAvailableを返します。
// 在庫不足になる場合はエラーを返します。
func (q *QuantityAvailable) Decrement(amount int) (*QuantityAvailable, error) {
	if amount < 0 {
		return nil, fmt.Errorf("負の値で減少させることはできません")
	}
	if q.value < amount {
		return nil, fmt.Errorf("在庫が不足しています")
	}
	return NewQuantityAvailable(q.value - amount)
}

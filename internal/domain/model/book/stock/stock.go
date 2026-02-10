package stock

import (
	"ddd-hands-on-go/internal/domain/model/book/stock/quantity_available"
	"ddd-hands-on-go/internal/domain/model/book/stock/status"
	"ddd-hands-on-go/internal/domain/model/book/stock/stock_id"
)

// Stock は在庫を表すエンティティです。
type Stock struct {
	stockId           *stock_id.StockId
	quantityAvailable *quantity_available.QuantityAvailable
	status            *status.Status
}

// NewStock は新しいStockを生成します。
func NewStock(id *stock_id.StockId, q *quantity_available.QuantityAvailable, s *status.Status) *Stock {
	return &Stock{
		stockId:           id,
		quantityAvailable: q,
		status:            s,
	}
}

// Create は初期状態のStockを生成します。
// IDの生成、初期在庫数(0)、ステータス(在庫切れ)の設定を行います。
func Create(idStr string) (*Stock, error) {
	id, err := stock_id.NewStockId(idStr)
	if err != nil {
		return nil, err
	}
	q, err := quantity_available.NewQuantityAvailable(0)
	if err != nil {
		return nil, err
	}
	s := status.NewStatus(status.OutOfStock)

	return NewStock(id, q, s), nil
}

// Reconstruct はDBなどから復元する際に使用するファクトリ関数です。
func Reconstruct(id *stock_id.StockId, q *quantity_available.QuantityAvailable, s *status.Status) *Stock {
	return NewStock(id, q, s)
}

// IncreaseQuantity は在庫数を増加させます。
// 在庫数が増加するとステータスが「在庫あり」になる可能性があります。
func (s *Stock) IncreaseQuantity(amount int) error {
	newQ, err := s.quantityAvailable.Increment(amount)
	if err != nil {
		return err
	}
	s.quantityAvailable = newQ

	if s.quantityAvailable.Value() > 0 {
		s.status = status.NewStatus(status.InStock)
	}
	return nil
}

// DecreaseQuantity は在庫数を減少させます。
// 在庫数が0になるとステータスが「在庫切れ」になります。
func (s *Stock) DecreaseQuantity(amount int) error {
	newQ, err := s.quantityAvailable.Decrement(amount)
	if err != nil {
		return err
	}
	s.quantityAvailable = newQ

	if s.quantityAvailable.Value() == 0 {
		s.status = status.NewStatus(status.OutOfStock)
	}
	return nil
}

// StockId は在庫IDを返します。
func (s *Stock) StockId() *stock_id.StockId {
	return s.stockId
}

// QuantityAvailable は在庫数を返します。
func (s *Stock) QuantityAvailable() *quantity_available.QuantityAvailable {
	return s.quantityAvailable
}

// Status は在庫ステータスを返します。
func (s *Stock) Status() *status.Status {
	return s.status
}

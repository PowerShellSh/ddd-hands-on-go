package domain_test

import (
	"ddd-hands-on-go/internal/domain/model/book"
	"ddd-hands-on-go/internal/domain/model/book/price"
	"ddd-hands-on-go/internal/domain/model/book/stock/status"
	"testing"
)

func TestBook(t *testing.T) {
	// 1. Value Objectの生成
	bookId, err := book.NewBookId("978-4-00-111111-1")
	if err != nil {
		t.Fatalf("BookIdの生成に失敗しました: %v", err)
	}

	title, err := book.NewTitle("Test Book")
	if err != nil {
		t.Fatalf("Titleの生成に失敗しました: %v", err)
	}

	p, err := price.NewPrice(1500, price.JPY)
	if err != nil {
		t.Fatalf("Priceの生成に失敗しました: %v", err)
	}

	// 2. Bookエンティティの生成
	b, err := book.NewBook(bookId, title, p)
	if err != nil {
		t.Fatalf("Bookの生成に失敗しました: %v", err)
	}

	// 3. 初期状態の検証
	if b.BookId().Value() != "978-4-00-111111-1" {
		t.Errorf("期待するBookId: 978-4-00-111111-1, 実際: %s", b.BookId().Value())
	}

	// 初期在庫は0、ステータスは在庫切れであるはず
	if b.Stock().QuantityAvailable().Value() != 0 {
		t.Errorf("期待する初期在庫: 0, 実際: %d", b.Stock().QuantityAvailable().Value())
	}
	if !b.Stock().Status().IsOutOfStock() {
		t.Errorf("期待する初期ステータス: OutOfStock")
	}

	// 4. 在庫の増加テスト
	if err := b.IncreaseStock(10); err != nil {
		t.Fatalf("在庫増加に失敗しました: %v", err)
	}

	if b.Stock().QuantityAvailable().Value() != 10 {
		t.Errorf("期待する在庫: 10, 実際: %d", b.Stock().QuantityAvailable().Value())
	}
	if b.Stock().Status().Value() != status.InStock {
		t.Errorf("期待するステータス: InStock")
	}

	// 5. 在庫の減少テスト
	if err := b.DecreaseStock(5); err != nil {
		t.Fatalf("在庫減少に失敗しました: %v", err)
	}
	if b.Stock().QuantityAvailable().Value() != 5 {
		t.Errorf("期待する在庫: 5, 実際: %d", b.Stock().QuantityAvailable().Value())
	}
}

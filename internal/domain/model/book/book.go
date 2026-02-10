package book

import (
	"ddd-hands-on-go/internal/domain/model/book/price"
	"ddd-hands-on-go/internal/domain/model/book/stock"
	"ddd-hands-on-go/internal/domain/shared"
	"fmt"
	"time"
)

// Book は書籍を表す集約ルートです。
type Book struct {
	bookId *BookId
	title  *Title
	price  *price.Price
	stock  *stock.Stock
	events []shared.DomainEvent // ドメインイベントを保持
}

// NewBook は新しいBookを生成します。
// 書籍生成時に、初期在庫（0、在庫切れ）も同時に生成されます。
func NewBook(bookId *BookId, title *Title, price *price.Price) (*Book, error) {
	// 新しい書籍を作成する際、在庫は初期状態で空/作成されます
	// 簡素化のため、StockIdにはBookIdと同じ値を使用します (1:1の関係)
	newStock, err := stock.Create(bookId.Value())
	if err != nil {
		return nil, fmt.Errorf("初期在庫の生成に失敗しました: %w", err)
	}

	book := &Book{
		bookId: bookId,
		title:  title,
		price:  price,
		stock:  newStock,
		events: make([]shared.DomainEvent, 0),
	}

	// BookCreatedイベントを記録
	book.AddEvent(&BookCreated{
		BookId:     bookId.Value(),
		Title:      title.Value(),
		OccurredAt: time.Now(),
	})

	return book, nil
}

// ReconstructBook はDBなどから復元する際に使用するファクトリ関数です。
func ReconstructBook(bookId *BookId, title *Title, price *price.Price, stock *stock.Stock) *Book {
	return &Book{
		bookId: bookId,
		title:  title,
		price:  price,
		stock:  stock,
		events: make([]shared.DomainEvent, 0),
	}
}

// ChangeTitle は書籍のタイトルを変更します。
func (b *Book) ChangeTitle(newTitle *Title) {
	b.title = newTitle
}

// ChangePrice は書籍の価格を変更します。
func (b *Book) ChangePrice(newPrice *price.Price) {
	b.price = newPrice
}

// IncreaseStock は在庫数を増加させます。
func (b *Book) IncreaseStock(amount int) error {
	return b.stock.IncreaseQuantity(amount)
}

// DecreaseStock は在庫数を減少させます。
func (b *Book) DecreaseStock(amount int) error {
	return b.stock.DecreaseQuantity(amount)
}

// IsSaleable は書籍が販売可能かどうかを判定します。
// 在庫があり、かつステータスが「在庫切れ」でない場合に販売可能とみなします。
func (b *Book) IsSaleable() bool {
	return b.stock.QuantityAvailable().Value() > 0 && !b.stock.Status().IsOutOfStock()
}

// BookId は書籍IDを返します。
func (b *Book) BookId() *BookId {
	return b.bookId
}

// Title は書籍タイトルを返します。
func (b *Book) Title() *Title {
	return b.title
}

// Price は書籍価格を返します。
func (b *Book) Price() *price.Price {
	return b.price
}

// Stock は在庫エンティティを返します。
func (b *Book) Stock() *stock.Stock {
	return b.stock
}

// AddEvent はドメインイベントを追加します。
func (b *Book) AddEvent(event shared.DomainEvent) {
	b.events = append(b.events, event)
}

// PullEvents は記録された全てのドメインイベントを取得し、内部リストをクリアします。
func (b *Book) PullEvents() []shared.DomainEvent {
	events := b.events
	b.events = make([]shared.DomainEvent, 0)
	return events
}

// BookCreated は書籍作成イベントです。
type BookCreated struct {
	BookId     string
	Title      string
	OccurredAt time.Time
}

func (e *BookCreated) EventName() string {
	return "BookCreated"
}

func (e *BookCreated) OccurredOn() time.Time {
	return e.OccurredAt
}

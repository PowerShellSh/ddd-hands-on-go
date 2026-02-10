package application_test

import (
	"context"
	"ddd-hands-on-go/internal/application/book"
	domain_book "ddd-hands-on-go/internal/domain/model/book"
	"ddd-hands-on-go/internal/domain/service"
	"testing"
)

// Mock Repositories
type mockBookRepository struct {
	books map[string]*domain_book.Book
}

func (m *mockBookRepository) Save(ctx context.Context, b *domain_book.Book) error {
	m.books[b.BookId().Value()] = b
	return nil
}

func (m *mockBookRepository) Find(ctx context.Context, bookId *domain_book.BookId) (*domain_book.Book, error) {
	if b, ok := m.books[bookId.Value()]; ok {
		return b, nil
	}
	return nil, nil // Not found
}

func (m *mockBookRepository) Delete(ctx context.Context, bookId *domain_book.BookId) error {
	delete(m.books, bookId.Value())
	return nil
}

type mockTransactionManager struct{}

func (m *mockTransactionManager) Begin(ctx context.Context, f func(ctx context.Context) error) error {
	return f(ctx)
}

func TestRegisterBookApplicationService(t *testing.T) {
	repo := &mockBookRepository{books: make(map[string]*domain_book.Book)}
	txManager := &mockTransactionManager{}
	dupSvc := service.NewISBNDuplicationCheckDomainService(repo)
	appSvc := book.NewRegisterBookApplicationService(repo, txManager, dupSvc)

	cmd := book.RegisterBookCommand{
		ISBN:        "978-4-00-111111-1",
		Title:       "Test Book",
		PriceAmount: 1500,
	}

	// 1. 正常系: 書籍登録成功
	if err := appSvc.Execute(context.Background(), cmd); err != nil {
		t.Fatalf("書籍登録に失敗しました: %v", err)
	}

	// 検証: リポジトリに保存されているか
	savedBook, err := repo.Find(context.Background(), mustBookId("978-4-00-111111-1"))
	if err != nil {
		t.Fatalf("書籍検索に失敗しました: %v", err)
	}
	if savedBook == nil {
		t.Fatalf("リポジトリに書籍が見つかりません")
	}
	if savedBook.Title().Value() != "Test Book" {
		t.Errorf("期待するタイトル: 'Test Book', 実際: %s", savedBook.Title().Value())
	}

	// 2. 異常系: 同一ISBNの重複登録 (ドメインサービスによる検証)
	if err := appSvc.Execute(context.Background(), cmd); err == nil {
		t.Errorf("重複ISBNのエラーが発生すべきですが、nilが返されました")
	}
}

func mustBookId(v string) *domain_book.BookId {
	id, _ := domain_book.NewBookId(v)
	return id
}

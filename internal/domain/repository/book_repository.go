package repository

import (
	"context"
	"ddd-hands-on-go/internal/domain/model/book"
)

// BookRepository は書籍エンティティの永続化を担当するリポジトリインターフェースです。
type BookRepository interface {
	// Save は書籍を保存します。
	Save(ctx context.Context, book *book.Book) error
	// Find は指定されたIDの書籍を検索して返します。
	Find(ctx context.Context, bookId *book.BookId) (*book.Book, error)
	// Delete は指定されたIDの書籍を削除します。
	Delete(ctx context.Context, bookId *book.BookId) error
}

package book

import (
	"context"
	"ddd-hands-on-go/internal/domain/model/book"
	"ddd-hands-on-go/internal/domain/repository"
)

// GetBookApplicationService は書籍情報取得ユースケースを実装するアプリケーションサービスです。
type GetBookApplicationService struct {
	bookRepository repository.BookRepository
}

// NewGetBookApplicationService は新しいGetBookApplicationServiceを生成します。
func NewGetBookApplicationService(bookRepository repository.BookRepository) *GetBookApplicationService {
	return &GetBookApplicationService{bookRepository: bookRepository}
}

// BookDTO は書籍情報のデータ転送オブジェクトです。
type BookDTO struct {
	ISBN              string  `json:"isbn"`
	Title             string  `json:"title"`
	PriceAmount       float64 `json:"price_amount"`
	QuantityAvailable int     `json:"quantity_available"`
	Status            string  `json:"status"`
}

// Execute は指定されたISBNの書籍情報を取得します。
func (s *GetBookApplicationService) Execute(ctx context.Context, isbn string) (*BookDTO, error) {
	bookId, err := book.NewBookId(isbn)
	if err != nil {
		return nil, err
	}

	foundBook, err := s.bookRepository.Find(ctx, bookId)
	if err != nil {
		return nil, err
	}
	if foundBook == nil {
		return nil, nil
	}

	return &BookDTO{
		ISBN:              foundBook.BookId().Value(),
		Title:             foundBook.Title().Value(),
		PriceAmount:       foundBook.Price().Amount(),
		QuantityAvailable: foundBook.Stock().QuantityAvailable().Value(),
		Status:            foundBook.Stock().Status().Value().String(),
	}, nil
}

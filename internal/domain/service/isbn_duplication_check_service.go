package service

import (
	"context"
	"ddd-hands-on-go/internal/domain/model/book"
	"ddd-hands-on-go/internal/domain/repository"
)

// ISBNDuplicationCheckDomainService はISBNの重複チェックを行うドメインサービスです。
type ISBNDuplicationCheckDomainService struct {
	bookRepository repository.BookRepository
}

// NewISBNDuplicationCheckDomainService は新しいISBNDuplicationCheckDomainServiceを生成します。
func NewISBNDuplicationCheckDomainService(bookRepository repository.BookRepository) *ISBNDuplicationCheckDomainService {
	return &ISBNDuplicationCheckDomainService{bookRepository: bookRepository}
}

// Execute はISBNが既に使用されているかどうかをチェックします。
func (s *ISBNDuplicationCheckDomainService) Execute(ctx context.Context, isbn *book.BookId) (bool, error) {
	foundBook, err := s.bookRepository.Find(ctx, isbn)
	if err != nil {
		// "見つからない" エラーであれば重複していないとみなす
		// エラーハンドリングはリポジトリの実装依存だが、ここでは簡易的にnilチェックを行う
		// 本来は errors.Is(err, repository.ErrNotFound) などで判定すべき
		return false, err
	}
	return foundBook != nil, nil
}

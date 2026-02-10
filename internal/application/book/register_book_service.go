package book

import (
	"context"
	"ddd-hands-on-go/internal/domain/model/book"
	"ddd-hands-on-go/internal/domain/model/book/price"
	"ddd-hands-on-go/internal/domain/repository"
	"ddd-hands-on-go/internal/domain/service"
	"ddd-hands-on-go/internal/domain/shared"
	"fmt"
)

// RegisterBookCommand は書籍登録に必要なパラーメータを保持する構造体です。
type RegisterBookCommand struct {
	ISBN        string
	Title       string
	PriceAmount float64
}

// RegisterBookApplicationService は書籍登録ユースケースを実装するアプリケーションサービスです。
type RegisterBookApplicationService struct {
	bookRepository          repository.BookRepository
	transactionManager      shared.TransactionManager
	duplicationCheckService *service.ISBNDuplicationCheckDomainService
	eventPublisher          shared.DomainEventPublisher
}

// NewRegisterBookApplicationService は新しいRegisterBookApplicationServiceを生成します。
func NewRegisterBookApplicationService(
	bookRepo repository.BookRepository,
	txManager shared.TransactionManager,
	dupCheck *service.ISBNDuplicationCheckDomainService,
	eventPublisher shared.DomainEventPublisher,
) *RegisterBookApplicationService {
	return &RegisterBookApplicationService{
		bookRepository:          bookRepo,
		transactionManager:      txManager,
		duplicationCheckService: dupCheck,
		eventPublisher:          eventPublisher,
	}
}

// Execute は書籍登録処理を実行します。
func (s *RegisterBookApplicationService) Execute(ctx context.Context, cmd RegisterBookCommand) error {
	return s.transactionManager.Begin(ctx, func(ctx context.Context) error {
		// 1. Value Objectの生成と検証
		isbn, err := book.NewBookId(cmd.ISBN)
		if err != nil {
			return err
		}
		title, err := book.NewTitle(cmd.Title)
		if err != nil {
			return err
		}
		// 通貨は現在JPY固定
		price, err := price.NewPrice(cmd.PriceAmount, price.JPY)
		if err != nil {
			return err
		}

		// 2. ISBN重複チェック
		isDuplicate, err := s.duplicationCheckService.Execute(ctx, isbn)
		if err != nil {
			return err
		}
		if isDuplicate {
			return fmt.Errorf("このISBNの書籍は既に存在します: %s", cmd.ISBN)
		}

		// 3. Bookエンティティの生成
		newBook, err := book.NewBook(isbn, title, price)
		if err != nil {
			return err
		}

		// 4. リポジトリへの保存
		if err := s.bookRepository.Save(ctx, newBook); err != nil {
			return err
		}

		// 5. ドメインイベントの発行
		// トランザクションコミット前に発行するか後に発行するかは要件次第だが、
		// ここではシンプルに処理成功の一部として発行する
		// (厳密にはTxコミット後に発行したいケースも多い)
		for _, event := range newBook.PullEvents() {
			s.eventPublisher.Publish(event)
		}

		return nil
	})
}

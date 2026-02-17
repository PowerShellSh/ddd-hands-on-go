//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	"ddd-hands-on-go/cmd/api/handler"
	"ddd-hands-on-go/internal/application/book"
	"ddd-hands-on-go/internal/domain/service"
	"ddd-hands-on-go/internal/infrastructure/event"
	"ddd-hands-on-go/internal/infrastructure/postgres"
	"ddd-hands-on-go/internal/infrastructure/subscriber"

	"github.com/google/wire"
)

// InitializeApp はアプリケーションの依存関係を解決し、ハンドラーを返します。
func InitializeApp(db *sql.DB) (*handler.BookHandler, error) {
	wire.Build(
		// リポジトリ
		postgres.NewPostgresBookRepository,
		wire.Bind(new(repository.BookRepository), new(*postgres.PostgresBookRepository)),

		// トランザクションマネージャー
		postgres.NewPostgresTransactionManager,
		wire.Bind(new(shared.TransactionManager), new(*postgres.PostgresTransactionManager)),

		// イベント配信
		event.NewEventEmitter,
		wire.Bind(new(shared.DomainEventPublisher), new(*event.EventEmitter)),

		// サブスクライバー
		subscriber.NewLogSubscriber,

		// ドメインサービス
		service.NewISBNDuplicationCheckDomainService,

		// アプリケーションサービス
		book.NewRegisterBookApplicationService,
		book.NewGetBookApplicationService,

		// ハンドラー
		handler.NewBookHandler,
	)
	return nil, nil
}

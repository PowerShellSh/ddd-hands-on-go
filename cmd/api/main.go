package main

import (
	"log"
	"net/http"
	"os"

	"ddd-hands-on-go/cmd/api/handler"
	"ddd-hands-on-go/internal/application/book"
	"ddd-hands-on-go/internal/domain/service"
	"ddd-hands-on-go/internal/infrastructure/event"
	"ddd-hands-on-go/internal/infrastructure/postgres"
	"ddd-hands-on-go/internal/infrastructure/subscriber"
)

func main() {
	// 環境変数の取得 (デフォルト値を設定)
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5433")
	dbUser := getEnv("DB_USER", "ddd")
	dbPass := getEnv("DB_PASS", "ddd")
	dbName := getEnv("DB_NAME", "ddd_hands_on")

	// 1. インフラストラクチャ層の初期化
	db, err := postgres.NewDB(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}
	defer db.Close()

	// 依存関係の注入
	bookRepo := postgres.NewPostgresBookRepository(db)
	txManager := postgres.NewPostgresTransactionManager(db)
	eventEmitter := event.NewEventEmitter()

	// サブスクライバーの登録
	logSubscriber := subscriber.NewLogSubscriber()
	eventEmitter.Subscribe("BookCreated", logSubscriber.Subscribe)

	// 2. ドメインサービスの初期化
	isbnDupCheckService := service.NewISBNDuplicationCheckDomainService(bookRepo)

	// 3. アプリケーションサービスの初期化
	registerBookService := book.NewRegisterBookApplicationService(bookRepo, txManager, isbnDupCheckService, eventEmitter)
	getBookService := book.NewGetBookApplicationService(bookRepo)

	// 4. ハンドラーの初期化
	bookHandler := handler.NewBookHandler(registerBookService, getBookService)

	// ルーティングの設定
	mux := http.NewServeMux()
	mux.HandleFunc("POST /books", bookHandler.RegisterBook)
	mux.HandleFunc("GET /books/{isbn}", bookHandler.GetBook)

	// サーバーの起動
	port := "8080"
	log.Printf("サーバーをポート %s で起動しています...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("サーバー起動エラー: %v", err)
	}
}

// getEnv は環境変数を取得し、設定されていない場合はデフォルト値を返します。
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

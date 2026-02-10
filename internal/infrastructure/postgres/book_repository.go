package postgres

import (
	"context"
	"database/sql"
	"ddd-hands-on-go/internal/domain/model/book"
	"ddd-hands-on-go/internal/domain/model/book/price"
	"ddd-hands-on-go/internal/domain/model/book/stock"
	"ddd-hands-on-go/internal/domain/model/book/stock/quantity_available"
	"ddd-hands-on-go/internal/domain/model/book/stock/status"
	"ddd-hands-on-go/internal/domain/model/book/stock/stock_id"
	"fmt"
)

// PostgresBookRepository はPostgreSQLを使用したBookRepositoryの実装です。
type PostgresBookRepository struct {
	db *sql.DB
}

// NewPostgresBookRepository は新しいPostgresBookRepositoryを生成します。
func NewPostgresBookRepository(db *sql.DB) *PostgresBookRepository {
	return &PostgresBookRepository{db: db}
}

// Save は書籍情報を保存(作成または更新)します。
func (r *PostgresBookRepository) Save(ctx context.Context, b *book.Book) error {
	var executor interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	}

	tx := GetTx(ctx)
	if tx != nil {
		executor = tx
	} else {
		executor = r.db
	}

	// Bookの保存 (Upsert)
	queryBook := `
		INSERT INTO "Book" ("bookId", "title", "priceAmount")
		VALUES ($1, $2, $3)
		ON CONFLICT ("bookId") DO UPDATE
		SET "title" = $2, "priceAmount" = $3
	`
	_, err := executor.ExecContext(ctx, queryBook,
		b.BookId().Value(),
		b.Title().Value(),
		b.Price().Amount(),
	)
	if err != nil {
		return fmt.Errorf("書籍の保存に失敗しました: %w", err)
	}

	// Stockの保存 (Upsert)
	queryStock := `
		INSERT INTO "Stock" ("stockId", "bookId", "quantityAvailable", "status")
		VALUES ($1, $2, $3, $4)
		ON CONFLICT ("stockId") DO UPDATE
		SET "quantityAvailable" = $3, "status" = $4
	`
	// Enumのマッピング
	statusStr := b.Stock().Status().Value().String()

	_, err = executor.ExecContext(ctx, queryStock,
		b.Stock().StockId().Value(),
		b.BookId().Value(),
		b.Stock().QuantityAvailable().Value(),
		statusStr,
	)
	if err != nil {
		return fmt.Errorf("在庫の保存に失敗しました: %w", err)
	}

	return nil
}

// Find は指定されたIDの書籍を検索します。
func (r *PostgresBookRepository) Find(ctx context.Context, bookId *book.BookId) (*book.Book, error) {
	query := `
		SELECT
			b."title",
			b."priceAmount",
			s."stockId",
			s."quantityAvailable",
			s."status"
		FROM "Book" b
		LEFT JOIN "Stock" s ON b."bookId" = s."bookId"
		WHERE b."bookId" = $1
	`

	row := r.db.QueryRowContext(ctx, query, bookId.Value())

	var titleStr string
	var priceAmount float64
	var stockIdStr string
	var quantityAvailableInt int
	var statusStr string

	if err := row.Scan(&titleStr, &priceAmount, &stockIdStr, &quantityAvailableInt, &statusStr); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 見つからない
		}
		return nil, fmt.Errorf("書籍の検索に失敗しました: %w", err)
	}

	title, err := book.NewTitle(titleStr)
	if err != nil {
		return nil, fmt.Errorf("データ整合性エラー: DB内のタイトルが不正です: %w", err)
	}

	price, err := price.NewPrice(priceAmount, price.JPY)
	if err != nil {
		return nil, fmt.Errorf("データ整合性エラー: DB内の価格が不正です: %w", err)
	}

	// Stockの再構築
	sId, err := stock_id.NewStockId(stockIdStr)
	if err != nil {
		return nil, fmt.Errorf("データ整合性エラー: DB内の在庫IDが不正です: %w", err)
	}

	q, err := quantity_available.NewQuantityAvailable(quantityAvailableInt)
	if err != nil {
		return nil, fmt.Errorf("データ整合性エラー: DB内の在庫数が不正です: %w", err)
	}

	stEnum := status.ToStatusEnum(statusStr)
	st := status.NewStatus(stEnum)

	stk := stock.Reconstruct(sId, q, st)

	return book.ReconstructBook(bookId, title, price, stk), nil
}

// Delete は書籍を削除します。
func (r *PostgresBookRepository) Delete(ctx context.Context, bookId *book.BookId) error {
	var executor interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	}

	tx := GetTx(ctx)
	if tx != nil {
		executor = tx
	} else {
		executor = r.db
	}

	query := `DELETE FROM "Book" WHERE "bookId" = $1`
	_, err := executor.ExecContext(ctx, query, bookId.Value())
	if err != nil {
		return fmt.Errorf("書籍の削除に失敗しました: %w", err)
	}

	return nil
}

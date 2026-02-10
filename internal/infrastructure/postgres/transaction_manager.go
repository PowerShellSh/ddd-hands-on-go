package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type key int

const (
	txKey key = iota
)

// PostgresTransactionManager はPostgreSQL用のトランザクションマネージャー実装です。
type PostgresTransactionManager struct {
	db *sql.DB
}

// NewPostgresTransactionManager は新しいPostgresTransactionManagerを生成します。
func NewPostgresTransactionManager(db *sql.DB) *PostgresTransactionManager {
	return &PostgresTransactionManager{db: db}
}

// Begin はトランザクションを開始します。
func (tm *PostgresTransactionManager) Begin(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("トランザクションの開始に失敗しました: %w", err)
	}

	// トランザクションをコンテキストに格納
	ctx = context.WithValue(ctx, txKey, tx)

	if err := f(ctx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("処理エラー: %w, ロールバックエラー: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("トランザクションのコミットに失敗しました: %w", err)
	}

	return nil
}

// GetTx はコンテキストからトランザクションを取得します。存在しない場合はnilを返します。
func GetTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return nil
}

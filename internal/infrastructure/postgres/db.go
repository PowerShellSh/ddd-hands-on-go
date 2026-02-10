package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// NewDB はPostgreSQLデータベースへの接続を確立します。
func NewDB(host, port, user, password, dbname string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("データベース接続のオープンに失敗しました: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("データベースへのPingに失敗しました: %w", err)
	}

	return db, nil
}

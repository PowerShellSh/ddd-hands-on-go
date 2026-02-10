package shared

import "context"

// TransactionManager はトランザクション制御を行うインターフェースです。
type TransactionManager interface {
	// Begin はトランザクションを開始し、関数内で実行される処理をアトミックに実行します。
	Begin(ctx context.Context, f func(ctx context.Context) error) error
}

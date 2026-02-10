package subscriber

import (
	"ddd-hands-on-go/internal/domain/model/book"
	"ddd-hands-on-go/internal/domain/shared"
	"log"
)

// LogSubscriber はイベントをログ出力するサブスクライバーです。
type LogSubscriber struct{}

// NewLogSubscriber は新しいLogSubscriberを生成します。
func NewLogSubscriber() *LogSubscriber {
	return &LogSubscriber{}
}

// Subscribe は指定されたイベントを購読します。
// ここではBookCreatedイベントを決め打ちで処理していますが、
// 汎用的にするためにインターフェースで受けて型アサーションすることも可能です。
func (s *LogSubscriber) Subscribe(event shared.DomainEvent) {
	if e, ok := event.(*book.BookCreated); ok {
		log.Printf("[LogSubscriber] 書籍が作成されました: ID=%s, Title=%s, At=%s",
			e.BookId, e.Title, e.OccurredAt)
	}
}

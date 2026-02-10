// DomainEvent はドメイン内で発生したイベントを表すインターフェースです。
package shared

import "time"

// DomainEvent は全てのドメインイベントが実装すべきインターフェースです。
type DomainEvent interface {
	EventName() string
	OccurredOn() time.Time
}

// DomainEventPublisher はドメインイベントを発行するためのインターフェースです。
type DomainEventPublisher interface {
	Publish(event DomainEvent)
}

// DomainEventSubscriber はドメインイベントを購読するためのインターフェースです。
type DomainEventSubscriber interface {
	Subscribe(eventName string, callback func(event DomainEvent))
}

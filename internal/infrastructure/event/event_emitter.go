package event

import (
	"ddd-hands-on-go/internal/domain/shared"
	"sync"
)

// EventEmitter はメモリ上でイベントの発行と購読を管理する構造体です。
type EventEmitter struct {
	subscribers map[string][]func(event shared.DomainEvent)
	mu          sync.RWMutex
}

// NewEventEmitter は新しいEventEmitterを生成します。
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		subscribers: make(map[string][]func(event shared.DomainEvent)),
	}
}

// Publish はイベントを発行し、登録されたリスナー（コールバック）を実行します。
func (e *EventEmitter) Publish(event shared.DomainEvent) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	callbacks, ok := e.subscribers[event.EventName()]
	if !ok {
		return
	}

	for _, callback := range callbacks {
		// 同期的に実行する（必要に応じてゴルーチンで非同期化も可能）
		callback(event)
	}
}

// Subscribe は指定されたイベント名に対してリスナーを登録します。
func (e *EventEmitter) Subscribe(eventName string, callback func(event shared.DomainEvent)) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.subscribers[eventName] = append(e.subscribers[eventName], callback)
}

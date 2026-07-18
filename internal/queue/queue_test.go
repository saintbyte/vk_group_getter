package queue

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestQueuePushAndStart(t *testing.T) {
	var received []uint
	var mu sync.Mutex

	q := New(10, func(id uint) {
		mu.Lock()
		received = append(received, id)
		mu.Unlock()
	})
	q.Start()

	q.Push(1)
	q.Push(2)
	q.Push(3)

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(received) != 3 {
		t.Fatalf("expected 3 items, got %d", len(received))
	}
	if received[0] != 1 || received[1] != 2 || received[2] != 3 {
		t.Errorf("unexpected order: %v", received)
	}
}

func TestQueueConcurrency(t *testing.T) {
	var count int64

	q := New(100, func(id uint) {
		atomic.AddInt64(&count, 1)
	})
	q.Start()

	for i := uint(0); i < 50; i++ {
		q.Push(i)
	}

	time.Sleep(100 * time.Millisecond)

	if atomic.LoadInt64(&count) != 50 {
		t.Errorf("expected 50 items processed, got %d", atomic.LoadInt64(&count))
	}
}

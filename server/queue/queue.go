package queue

import (
	"fmt"
	"sync"
)

type Queue[T any] struct {
	lock  sync.RWMutex
	queue []T
}

func CreateQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) Enqueue(data T) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue = append(q.queue, data)
}

func (q *Queue[T]) Dequeue() T {
	q.lock.Lock()
	defer q.lock.Unlock()
	var val T = q.queue[0]
	q.queue = q.queue[1:]
	return val
}

func (q *Queue[T]) Peek() any {
	if q.IsEmpty() {
		return ""
	}
	return q.queue[0]
}

func (q *Queue[T]) IsEmpty() bool {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return len(q.queue) <= 0
}

func (q *Queue[T]) Display() {
	for i := 0; i < len(q.queue); i++ {
		fmt.Println(q.queue[i])
	}
}

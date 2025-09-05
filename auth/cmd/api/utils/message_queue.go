package utils

import "sync"

type MessageQueue[T any] struct {
	mutex sync.Mutex
	cond  *sync.Cond
	queue []T
}

func NewMessageQueue[T any]() *MessageQueue[T] {
	messageQueue := new(MessageQueue[T])
	messageQueue.cond = sync.NewCond(&messageQueue.mutex)
	return messageQueue
}

func (queue *MessageQueue[T]) Enqueue(item T) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	queue.queue = append(queue.queue, item)
	queue.cond.Signal()
}

func (queue *MessageQueue[T]) Dequeue() T {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	for len(queue.queue) == 0 {
		queue.cond.Wait()
	}

	item := queue.queue[0]
	queue.queue = queue.queue[1:]

	return item
}

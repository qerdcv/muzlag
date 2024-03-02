package queue

import (
	"errors"
	"sync"
)

var ErrEmptyQueue = errors.New("empty queue")

type Item struct {
	Title     string
	OriginURL string
}

type Queue struct {
	mu sync.Mutex

	queue map[string][]Item
}

func New() *Queue {
	return &Queue{
		queue: map[string][]Item{},
	}
}

func (q *Queue) AddToQueue(qID string, i Item) {
	q.mu.Lock()

	defer q.mu.Unlock()

	if _, ok := q.queue[qID]; !ok {
		q.queue[qID] = make([]Item, 0)
	}

	q.queue[qID] = append(q.queue[qID], i)
}

func (q *Queue) PopQueue(qID string) (Item, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if queue, ok := q.queue[qID]; ok {
		if len(queue) == 0 {
			delete(q.queue, qID)
			return Item{}, ErrEmptyQueue
		}

		var i Item
		i, q.queue[qID] = q.queue[qID][0], q.queue[qID][1:]

		return i, nil
	}

	return Item{}, ErrEmptyQueue
}

func (q *Queue) Next(qID string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	queue, ok := q.queue[qID]
	return ok && len(queue) != 0
}

func (q *Queue) List(qID string) ([]Item, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if queue, ok := q.queue[qID]; ok && len(queue) != 0 {
		return queue, nil
	}

	return nil, ErrEmptyQueue
}

func (q *Queue) CleanQueue(qID string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	delete(q.queue, qID)
}

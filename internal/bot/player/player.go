package player

import (
	"errors"
	"sync"
)

var ErrEmptyQueue = errors.New("empty queue")

type Player struct {
	mu sync.Mutex

	musicQueue map[string][]string
}

func New() *Player {
	return &Player{
		musicQueue: map[string][]string{},
	}
}

func (p *Player) AddToQueue(qID, url string) {
	p.mu.Lock()

	defer p.mu.Unlock()

	if _, ok := p.musicQueue[qID]; !ok {
		p.musicQueue[qID] = make([]string, 0)
	}

	p.musicQueue[qID] = append(p.musicQueue[qID], url)
}

func (p *Player) PopQueue(qID string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if q, ok := p.musicQueue[qID]; ok {
		if len(q) == 0 {
			delete(p.musicQueue, qID)
			return "", ErrEmptyQueue
		}

		var sURL string
		sURL, p.musicQueue[qID] = p.musicQueue[qID][0], p.musicQueue[qID][1:]

		return sURL, nil
	}

	return "", ErrEmptyQueue
}

func (p *Player) Next(qID string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	q, ok := p.musicQueue[qID]
	return ok && len(q) != 0
}

func (p *Player) CleanQueue(qID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.musicQueue, qID)
}

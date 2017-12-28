package session

import (
	"log"
	"sync"
	"time"
)

type Storage interface {
	Check(string) bool
	Get(string) map[string]interface{}
	Expire()
	Add(string, map[string]interface{})
	Remove(string)
}

func NewInMemoryStorage(maxSize int, maxIdle, maxAge time.Duration) Storage {
	return &inMemoryStorage{
		maxIdle: maxIdle,
		maxAge:  maxAge,
		maxSize: maxSize,
		tokens:  make(map[string]storageEntry),
		mutex:   &sync.Mutex{},
	}
}

type inMemoryStorage struct {
	maxIdle time.Duration
	maxAge  time.Duration
	maxSize int
	tokens  map[string]storageEntry
	mutex   *sync.Mutex
}

func (s *inMemoryStorage) Check(token string) bool {
	entry, ok := s.tokens[token]
	if !ok {
		return false
	}

	entry.Touch()

	log.Printf("Touched entry %s, age %s", entry.Idle(), entry.Age())

	return true
}

func (s *inMemoryStorage) Get(token string) map[string]interface{} {
	entry, ok := s.tokens[token]
	if !ok {
		return nil
	}

	return entry.data
}

func (s *inMemoryStorage) Expire() {
	s.mutex.Lock()

	remove := []string{}
	for t, entry := range s.tokens {
		if entry.Age() > s.maxAge || entry.Idle() > s.maxIdle {
			remove = append(remove, t)
		}
	}

	for _, t := range remove {
		log.Printf("expiring token %s", t)
		delete(s.tokens, t)
	}

	s.mutex.Unlock()
}

func (s *inMemoryStorage) Add(token string, values map[string]interface{}) {
	s.Expire()
	s.mutex.Lock()

	s.tokens[token] = storageEntry{
		createdAt:  time.Now().UTC(),
		accessedAt: time.Now().UTC(),
		data:       values,
	}

	s.mutex.Unlock()
}

func (s *inMemoryStorage) Remove(token string) {
}

type storageEntry struct {
	createdAt  time.Time
	accessedAt time.Time
	data       map[string]interface{}
}

func (e storageEntry) Age() time.Duration {
	return time.Since(e.createdAt)
}

func (e storageEntry) Idle() time.Duration {
	return time.Since(e.accessedAt)
}

func (e *storageEntry) Touch() {
	e.accessedAt = time.Now()
}

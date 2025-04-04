package adapters

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type RoutesRepository interface {
	Gen() (string, error)
	Exists(code string) bool
}

type InMemoryRoutesRepository struct {
	mu    sync.RWMutex
	codes map[string]string
}

func NewInMemoryRoutesRepository() *InMemoryRoutesRepository {
	return &InMemoryRoutesRepository{
		codes: make(map[string]string),
	}
}

func (r *InMemoryRoutesRepository) Gen(fixedId string) (string, error) {
	for {
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			return "", err
		}

		// b = append(b, []byte(fixedId)...)

		code := hex.EncodeToString(b)

		if _, exists := r.Exists(code); !exists {
			r.mu.Lock()
			r.codes[code] = fixedId
			r.mu.Unlock()
			return code, nil
		}
	}
}

func (r *InMemoryRoutesRepository) Exists(code string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, exists := r.codes[code]
	return c, exists
}

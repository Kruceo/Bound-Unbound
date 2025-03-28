package adapters

import (
	"fmt"
	"sync"
)

type InMemoryResponseRepository struct {
	data    map[string]string
	mu      sync.Mutex
	channel chan string
}

func (r *InMemoryResponseRepository) Set(id string, data string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[id] = data
	r.channel <- id
	return nil
}

func (r *InMemoryResponseRepository) WaitForResponse(id string) error {
	for v := range r.channel {
		if v == id {
			return nil
		}
	}
	return nil
}

func (r *InMemoryResponseRepository) ReadResponse(id string) (string, error) {
	return r.data[id], nil
}

func (r *InMemoryResponseRepository) DeleteResponse(id string) error {
	return fmt.Errorf("not implemented")
}

func NewInMemoryResponseRepository() InMemoryResponseRepository {
	return InMemoryResponseRepository{data: make(map[string]string), mu: sync.Mutex{}, channel: make(chan string)}
}

// type ResponsesReporisory interface {
// 	Set(id string, data string) error
// 	WaitForResponse(id string) error
// 	ReadResponse(id string) (string, error)
// 	DeleteResponse(id string) error
// }

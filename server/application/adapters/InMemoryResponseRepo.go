package adapters

import (
	"fmt"
	"sync"
	"time"
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
	go func() {
		ticker := time.NewTimer(30 * time.Second)
		defer ticker.Stop()
		<-ticker.C
		fmt.Println("timeout for", id)
		r.channel <- "_TIMEOUT_" + id
	}()

	for v := range r.channel {
		if v == id {
			return nil
		}
		if v == "_TIMEOUT_"+id {
			return fmt.Errorf("timeout for response id %s", id)
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

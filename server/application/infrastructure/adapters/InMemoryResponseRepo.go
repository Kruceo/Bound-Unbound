package adapters

import (
	"fmt"
	"sync"
	"time"
)

type InMemoryResponseRepository struct {
	data     map[string]string
	mu       sync.Mutex
	channels map[string]chan string
}

func (r *InMemoryResponseRepository) Set(id string, data string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case r.channels[id] <- id:
		r.data[id] = data
	case <-time.After(2 * time.Second):
		return fmt.Errorf("timeout to set response: %s", id)
	}
	return nil
}

func (r *InMemoryResponseRepository) WaitForResponse(id string) error {
	r.mu.Lock()
	if _, exists := r.channels[id]; !exists {
		fmt.Println("criando canal no waitresponse")
		r.channels[id] = make(chan string, 2)
	}
	fmt.Println("passou daqui")
	ch := r.channels[id]
	r.mu.Unlock()

	go func() {
		time.Sleep(5 * time.Second)
		ch <- "_TIMEOUT_" + id
	}()

	for v := range ch {
		fmt.Println("channel kkk", v, id)
		if v == id {
			return nil
		}
		if v == "_TIMEOUT_"+id {
			fmt.Println("timeout for", id)
			return fmt.Errorf("timeout for response id %s", id)
		}
	}

	return nil
}

func (r *InMemoryResponseRepository) ReadResponse(id string) (string, error) {
	r.mu.Lock()
	d := r.data[id]
	delete(r.data, id)
	r.mu.Unlock()
	return d, nil
}

func (r *InMemoryResponseRepository) DeleteResponse(id string) error {
	return fmt.Errorf("not implemented")
}

func NewInMemoryResponseRepository() InMemoryResponseRepository {
	return InMemoryResponseRepository{data: make(map[string]string), mu: sync.Mutex{}, channels: make(map[string]chan string)}
}

// type ResponsesReporisory interface {
// 	Set(id string, data string) error
// 	WaitForResponse(id string) error
// 	ReadResponse(id string) (string, error)
// 	DeleteResponse(id string) error
// }

package adapters

import (
	"server2/application/entities"
	"sync"
)

type InMemoryNodeRepository struct {
	data map[string]entities.Node
	mu   sync.Mutex
}

func (r *InMemoryNodeRepository) Save(node entities.Node) (string, error) {
	// id := fmt.Sprintf("%x", rand.Uint32())
	id := node.Conn.RemoteAddr().String()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[id] = node
	return id, nil
}

func (r *InMemoryNodeRepository) Get(id string) *entities.Node {
	r.mu.Lock() // Pode ser necessário usar um mutex no Get também, dependendo do uso concorrente
	defer r.mu.Unlock()
	n := r.data[id]
	return &n
}

func (r *InMemoryNodeRepository) IDs() []string {
	ids := []string{}
	for id, _ := range r.data {
		ids = append(ids, id)
	}
	return ids
}

func (r *InMemoryNodeRepository) Delete(id string) error {
	delete(r.data, id)
	return nil
}

func NewInMemoryNodeRepository() InMemoryNodeRepository {
	return InMemoryNodeRepository{data: make(map[string]entities.Node), mu: sync.Mutex{}}
}

//  Save(node Node) (string, error)
// 	Get(id string) *Node
// 	IDs() []string

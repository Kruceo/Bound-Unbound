package adapters

import (
	"crypto/cipher"
	"fmt"
	"server2/application/entities"
	"sync"

	"github.com/gorilla/websocket"
)

type InMemoryNodeRepository struct {
	data map[string]entities.Node
	mu   sync.Mutex
}

func (r *InMemoryNodeRepository) Save(id, name string, conn *websocket.Conn, cipher *cipher.AEAD) (string, error) {
	// id := fmt.Sprintf("%x", rand.Uint32())
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[id] = *entities.NewNode(id, name, conn, cipher)
	return id, nil
}

func (r *InMemoryNodeRepository) Get(id string) *entities.Node {
	r.mu.Lock() // Pode ser necessário usar um mutex no Get também, dependendo do uso concorrente
	defer r.mu.Unlock()
	n, exists := r.data[id]
	if !exists {
		return nil
	}
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

func (r *InMemoryNodeRepository) FindOneByRemoteAddress(remoteAddr string) (*entities.Node, error) {
	for _, v := range r.data {
		if v.Conn.RemoteAddr().String() == remoteAddr {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("node not found: %s", remoteAddr)
}

func NewInMemoryNodeRepository() InMemoryNodeRepository {
	return InMemoryNodeRepository{data: make(map[string]entities.Node), mu: sync.Mutex{}}
}

//  Save(node Node) (string, error)
// 	Get(id string) *Node
// 	IDs() []string

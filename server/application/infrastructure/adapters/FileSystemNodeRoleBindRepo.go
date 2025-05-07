package adapters

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"server2/application/entities"

	"github.com/google/uuid"
)

type FileNodeRoleBindRepository struct {
	filePath string
	storage  map[string]*entities.RoleNodeBind
	mutex    sync.RWMutex
}

func NewFileNodeRoleBindRepository(filePath string) (*FileNodeRoleBindRepository, error) {
	repo := &FileNodeRoleBindRepository{
		filePath: filePath,
		storage:  make(map[string]*entities.RoleNodeBind),
	}

	err := repo.loadFromFile()
	if err != nil {

		return nil, err
	}
	return repo, nil
}

func (repo *FileNodeRoleBindRepository) loadFromFile() error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	file, err := os.Open(repo.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // arquivo ainda não existe, sem erro
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&repo.storage)
}

func (repo *FileNodeRoleBindRepository) saveToFile() error {
	fmt.Println("save")
	locked := repo.mutex.TryLock()
	if locked {
		defer repo.mutex.Unlock()
	}

	file, err := os.Create(repo.filePath)
	fmt.Println("save2")
	if err != nil {
		fmt.Printf("save7")
		return err
	}
	defer file.Close()

	fmt.Println("save3")
	encoder := json.NewEncoder(file)
	fmt.Println("save4")
	encoder.SetIndent("", "  ")
	fmt.Println("save5")
	return encoder.Encode(repo.storage)
}

func (repo *FileNodeRoleBindRepository) Save(id, nodeID, roleID string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	r, err := entities.NewRoleNodeBind(id, nodeID, roleID)
	if err != nil {
		return err
	}
	repo.storage[id] = r

	return repo.saveToFile()
}

func (repo *FileNodeRoleBindRepository) Get(id string) (*entities.RoleNodeBind, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	bind, exists := repo.storage[id]
	if !exists {
		return nil, errors.New("bind not found")
	}

	return bind, nil
}

func (repo *FileNodeRoleBindRepository) GetAll(limit int) ([]*entities.RoleNodeBind, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	list := []*entities.RoleNodeBind{}
	count := 0
	for _, v := range repo.storage {
		if count > limit {
			break
		}
		list = append(list, v)
		count++
	}
	return list, nil
}

func (repo *FileNodeRoleBindRepository) Delete(id string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if _, exists := repo.storage[id]; !exists {
		return errors.New("bind not found")
	}

	delete(repo.storage, id)
	return repo.saveToFile()
}

func (repo *FileNodeRoleBindRepository) Update(id, nodeID, roleID string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if _, exists := repo.storage[id]; !exists {
		return errors.New("bind not found")
	}

	repo.storage[id] = &entities.RoleNodeBind{
		ID:     id,
		NodeID: nodeID,
		RoleID: roleID,
	}

	return repo.saveToFile()
}

func (repo *FileNodeRoleBindRepository) NextID() (string, error) {
	return uuid.NewString(), nil
}

package adapters

import (
	"encoding/json"
	"errors"
	"os"
	"server2/application/entities"
	"server2/utils"
	"sync"

	"github.com/google/uuid"
)

type FileRoleRepository struct {
	filePath string
	mutex    sync.Mutex
	roles    map[string]*entities.Role
}

func NewFileRoleRepository(path string) *FileRoleRepository {
	repo := &FileRoleRepository{
		filePath: path,
		roles:    make(map[string]*entities.Role),
	}

	_, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	err = repo.load()
	if err != nil {
		panic(err)
	}

	admin, err := entities.NewRole("0", "Admin", utils.ADMIN_PERMS...)
	if err != nil {
		panic(err)
	}

	_, err = repo.CreateIfNotExists(admin)
	if err != nil {
		panic(err)
	}
	return repo
}

func (r *FileRoleRepository) load() error {
	file, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &r.roles)
}

func (r *FileRoleRepository) save() error {
	data, err := json.MarshalIndent(r.roles, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}

func (r *FileRoleRepository) Create(role *entities.Role) (string, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.roles[role.ID]; exists {
		return "", errors.New("role already exists")
	}

	r.roles[role.ID] = role
	return role.ID, r.save()
}

func (r *FileRoleRepository) Get(id string) (*entities.Role, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	role, ok := r.roles[id]
	if !ok {
		return nil, errors.New("role not found")
	}
	return role, nil
}

func (r *FileRoleRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, ok := r.roles[id]
	if !ok {
		return errors.New("role not found")
	}

	delete(r.roles, id)
	return r.save()
}

func (r *FileRoleRepository) Update(role *entities.Role) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.roles[role.ID]; !ok {
		return errors.New("role not found")
	}

	r.roles[role.ID] = role
	return r.save()
}

func (r *FileRoleRepository) GetAll(limit int) ([]*entities.Role, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	roles := make([]*entities.Role, 0, len(r.roles))
	count := 0
	for _, role := range r.roles {
		roles = append(roles, role)
		count++
		if limit > 0 && count >= limit {
			break
		}
	}
	return roles, nil
}

func (r *FileRoleRepository) SearchByName(name string, limit int) ([]*entities.Role, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var result []*entities.Role
	for _, role := range r.roles {
		if name == "" || containsIgnoreCase(role.Name, name) {
			result = append(result, role)
			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (r *FileRoleRepository) NextID() (string, error) {
	id := uuid.New()
	return id.String(), nil
}

func (r *FileRoleRepository) CreateIfNotExists(role *entities.Role) (bool, error) {
	if _, exists := r.roles[role.ID]; exists {
		return false, nil
	}
	_, err := r.Create(role)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *FileRoleRepository) Count() (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return len(r.roles), nil
}

func containsIgnoreCase(a, b string) bool {
	return len(a) >= len(b) && stringEqualFold(a[:len(b)], b)
}

// compara strings ignorando case (simplificado)
func stringEqualFold(a, b string) bool {
	return len(a) == len(b) && (a == b || toLower(a) == toLower(b))
}

func toLower(s string) string {
	result := []rune{}
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			r = r + 32
		}
		result = append(result, r)
	}
	return string(result)
}

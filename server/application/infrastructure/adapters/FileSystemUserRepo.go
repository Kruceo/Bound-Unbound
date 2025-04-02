package adapters

import (
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"server2/application/entities"
	"sync"
)

type FileSystemUserRepo struct {
	filePath string
	mutex    sync.Mutex
}

func NewFileSystemUserRepo(filePath string) *FileSystemUserRepo {
	return &FileSystemUserRepo{filePath: filePath}
}

func (f *FileSystemUserRepo) loadUsers() (map[string]*entities.User, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	file, err := os.Open(f.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]*entities.User), nil
		}
		return nil, err
	}
	defer file.Close()

	var users map[string]*entities.User
	if err := json.NewDecoder(file).Decode(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func (f *FileSystemUserRepo) saveUsers(users map[string]*entities.User) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	file, err := os.Create(f.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(users)
}

func (f *FileSystemUserRepo) Save(name, password string, role uint8, recoveryCode string) (string, error) {
	users, err := f.loadUsers()
	if err != nil {
		return "", err
	}

	id := name // Usando o nome como ID (ajuste se precisar)
	users[id] = &entities.User{Username: name, Password: password, Role: role, RecoveryCode: recoveryCode}

	if err := f.saveUsers(users); err != nil {
		return "", err
	}
	return id, nil
}

func (f *FileSystemUserRepo) Get(id string) (*entities.User, error) {
	users, err := f.loadUsers()
	if err != nil {
		return nil, err
	}

	user, exists := users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (f *FileSystemUserRepo) Delete(id string) error {
	users, err := f.loadUsers()
	if err != nil {
		return err
	}

	if _, exists := users[id]; !exists {
		return errors.New("user not found")
	}
	delete(users, id)

	return f.saveUsers(users)
}

func (f *FileSystemUserRepo) SearchByName(regex string) ([]*entities.User, error) {
	users, err := f.loadUsers()
	if err != nil {
		return nil, err
	}

	var result []*entities.User
	r := regexp.MustCompile(regex)
	for _, user := range users {
		// fmt.Println(user.Username, r, r.MatchString(user.Username))
		if r.MatchString(user.Username) {
			result = append(result, user)
		}
	}
	return result, nil
}

func (f *FileSystemUserRepo) SearchByRole(role uint8) ([]*entities.User, error) {
	users, err := f.loadUsers()
	if err != nil {
		return nil, err
	}

	var result []*entities.User
	for _, user := range users {
		if user.Role == role {
			result = append(result, user)
		}
	}
	return result, nil
}

func (f *FileSystemUserRepo) FindOneByName(regex string) (*entities.User, error) {
	users, err := f.SearchByName(regex)
	if err != nil || len(users) == 0 {
		return nil, errors.New("user not found")
	}
	return users[0], nil
}

func (f *FileSystemUserRepo) FindOneByRole(role uint8) (*entities.User, error) {
	users, err := f.SearchByRole(role)
	if err != nil || len(users) == 0 {
		return nil, errors.New("user not found")
	}
	return users[0], nil
}

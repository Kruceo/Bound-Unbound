package adapters

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (f *FileSystemUserRepo) Save(name, password string, roleID string, recoveryCode string) (string, error) {
	users, err := f.loadUsers()
	if err != nil {
		return "", err
	}

	id := name // in newer repositories use other things as id (uuid, sequential), this adapter is temp

	if _, exists := users[id]; exists {
		return "", fmt.Errorf("user id already exists: %s", id)
	}

	user, err := entities.NewUser(id, name, password, roleID, recoveryCode)
	if err != nil {
		return "", err
	}

	users[id] = &user

	if err := f.saveUsers(users); err != nil {
		return "", err
	}
	return id, nil
}

func (f *FileSystemUserRepo) Update(id, name, password string, roleID string, secretCodeHash string) error {
	users, err := f.loadUsers()
	if err != nil {
		return err
	}
	if _, exists := users[id]; !exists {
		return fmt.Errorf("user not found: %s", id)
	}
	users[id].Username = name
	users[id].SetPassword(password)
	users[id].RoleID = roleID
	users[id].RecoveryCode = secretCodeHash

	err = f.saveUsers(users)
	if err != nil {
		return err
	}

	return nil
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
		if r.MatchString(user.Username) {
			result = append(result, user)
		}
	}
	return result, nil
}

func (f *FileSystemUserRepo) SearchByRoleID(role string) ([]*entities.User, error) {
	users, err := f.loadUsers()
	if err != nil {
		return nil, err
	}

	var result []*entities.User
	for _, user := range users {
		if user.RoleID == fmt.Sprintf("%s", role) {
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

func (f *FileSystemUserRepo) FindOneByRoleID(role string) (*entities.User, error) {
	users, err := f.SearchByRoleID(role)
	if err != nil || len(users) == 0 {
		return nil, errors.New("user not found")
	}
	return users[0], nil
}

func (f *FileSystemUserRepo) Count() (int, error) {
	users, err := f.loadUsers()
	if err != nil {
		return 0, err
	}
	return len(users), nil
}

func (f *FileSystemUserRepo) CountByRoleID(role string) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (f *FileSystemUserRepo) CountByName(name string) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

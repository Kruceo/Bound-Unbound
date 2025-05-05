package usecases

import (
	"fmt"
	"server2/application/entities"
	"server2/application/infrastructure"
)

type RoleUseCase struct {
	repo infrastructure.RoleRepository
}

func NewRoleUseCase(repo infrastructure.RoleRepository) *RoleUseCase {
	return &RoleUseCase{
		repo: repo,
	}
}

// return createdId,err
func (u *RoleUseCase) Save(name string, permissions []string) (string, error) {
	id, err := u.repo.NextID()
	if err != nil {
		return "", err
	}

	r, err := entities.NewRole(id, name, permissions...)
	if err != nil {
		return "", err
	}
	return u.repo.Create(r)
}

func (u *RoleUseCase) Get(id string) (*entities.Role, error) {
	fmt.Println("get role", id)
	return u.repo.Get(id)
}

func (u *RoleUseCase) Delete(id string) error {
	return u.repo.Delete(id)
}

func (u *RoleUseCase) Update(role *entities.Role) error {
	return u.repo.Update(role)
}

func (u *RoleUseCase) GetAll(limit int) ([]*entities.Role, error) {
	return u.repo.GetAll(limit)
}

func (u *RoleUseCase) SearchByName(name string, limit int) ([]*entities.Role, error) {
	return u.repo.SearchByName(name, limit)
}

func (u *RoleUseCase) Count() (int, error) {
	return u.repo.Count()
}

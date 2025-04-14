package usecases

import (
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
func (u *RoleUseCase) Create(role *entities.Role) (string, error) {
	return u.repo.Create(role)
}

func (u *RoleUseCase) Get(id string) (*entities.Role, error) {
	return u.repo.Get(id)
}

func (u *RoleUseCase) Delete(id string) (*entities.Role, error) {
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

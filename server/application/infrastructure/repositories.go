package infrastructure

import "server2/application/entities"

type NodeRepository interface {
	Save(node entities.Node) (string, error)
	Get(id string) *entities.Node
	Delete(id string) error
	IDs() []string
}

type ResponsesReporisory interface {
	Set(id string, data string) error
	WaitForResponse(id string) error
	ReadResponse(id string) (string, error)
	DeleteResponse(id string) error
}

type UserRepository interface {
	Save(name, password string, roleID string, recoveryCode string) (string, error)
	Update(id, name, password string, roleID string, secretCode string) error
	Get(id string) (*entities.User, error)
	Delete(id string) error
	SearchByName(regex string) ([]*entities.User, error)
	SearchByRoleName(role string) ([]*entities.User, error)
	FindOneByName(regex string) (*entities.User, error)
	FindOneByRoleName(role string) (*entities.User, error)
	Count() (int, error)
	CountByRoleName(role string) (int, error)
	CountByName(regex string) (int, error)
}

type RoutesRepository interface {
	Gen(string) (string, error)
	Exists(string) (string, bool)
}

type RequestBlocker interface {
	IsBlocked(ip string) bool
	MarkAttempt(ip string)
}

type RoleRepository interface {
	Create(*entities.Role) (string, error)
	Get(id string) (*entities.Role, error)
	Delete(id string) (*entities.Role, error)
	Update(*entities.Role) error
	GetAll(limit int) ([]*entities.Role, error)
	SearchByName(name string, limit int) ([]*entities.Role, error)
	Count() (int, error)
}

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
	Save(name, password string, role uint8, recoveryCode string) (string, error)
	Update(id, name, password string, role uint8, secretCode string) error
	Get(id string) (*entities.User, error)
	Delete(id string) error
	SearchByName(regex string) ([]*entities.User, error)
	SearchByRole(role uint8) ([]*entities.User, error)
	FindOneByName(regex string) (*entities.User, error)
	FindOneByRole(role uint8) (*entities.User, error)
}

type RoutesRepository interface {
	Gen(string) (string, error)
	Exists(string) (string, bool)
}

type RequestBlocker interface {
	IsBlocked(ip string) bool
	MarkAttempt(ip string)
}

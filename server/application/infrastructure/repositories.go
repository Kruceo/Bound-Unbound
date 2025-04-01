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

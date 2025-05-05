package usecases

import (
	"fmt"
	"server2/application/entities"
	"server2/application/infrastructure"
)

type NodeRoleBindPersistenceUseCase struct {
	nrbRepo  infrastructure.NodeRoleBindRepository
	nodeRepo infrastructure.NodeRepository
	roleRepo infrastructure.RoleRepository
}

func NewNodeRoleBindPersistenceUseCase(
	nrbRepo infrastructure.NodeRoleBindRepository,
	nodeRepo infrastructure.NodeRepository,
	roleRepo infrastructure.RoleRepository) *NodeRoleBindPersistenceUseCase {
	return &NodeRoleBindPersistenceUseCase{
		nrbRepo:  nrbRepo,
		nodeRepo: nodeRepo,
		roleRepo: roleRepo,
	}
}

func (nr *NodeRoleBindPersistenceUseCase) Bind(nodeID string, roleID string) (string, error) {
	_, err := nr.nodeRepo.Get(nodeID)
	if err != nil {
		return "", err
	}
	_, err = nr.roleRepo.Get(roleID)
	if err != nil {
		return "", err
	}

	newID, err := nr.nrbRepo.NextID()
	if err != nil {
		return "", err
	}
	return newID, nr.nrbRepo.Save(newID, nodeID, roleID)
}

func (nr *NodeRoleBindPersistenceUseCase) Save(nodeID string, roleID string) (string, error) {
	fmt.Println("saving bind", nodeID, roleID)
	id, err := nr.nrbRepo.NextID()
	if err != nil {
		return "", err
	}
	err = nr.nrbRepo.Save(id, nodeID, roleID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (nr *NodeRoleBindPersistenceUseCase) Get(id string) (*entities.RoleNodeBind, error) {
	return nr.nrbRepo.Get(id)
}

type GetAllUseCaseResponse struct {
	ID   string
	Node entities.Node
	Role entities.Role
}

func (nr *NodeRoleBindPersistenceUseCase) GetAllWithIncluded(limit int) ([]GetAllUseCaseResponse, error) {
	recovered, err := nr.nrbRepo.GetAll(limit)
	if err != nil { // handle unhandled error in this function
		return nil, err
	}
	res := make([]GetAllUseCaseResponse, 0, len(recovered))

	for _, v := range recovered {
		node, err := nr.nodeRepo.Get(v.NodeID)
		if err != nil {
			return nil, err
		}
		role, err := nr.roleRepo.Get(v.RoleID)
		if err != nil {
			return nil, err
		}
		res = append(res, GetAllUseCaseResponse{ID: v.ID, Node: *node, Role: *role})
	}
	return res, nil
}

func (nr *NodeRoleBindPersistenceUseCase) GetNodesForRole(roleID string) ([]*entities.Node, error) {
	fmt.Println("getting nodes for role", roleID)
	role, err := nr.roleRepo.Get(roleID)
	if err != nil {
		return nil, err
	}

	if role.HasPerm("admin") {
		ids := nr.nodeRepo.IDs()
		nodes := []*entities.Node{}
		for _, v := range ids {
			node, err := nr.nodeRepo.Get(v)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		}
		return nodes, nil
	}
	binds, err := nr.nrbRepo.GetAll(9999999)
	if err != nil {
		return nil, err
	}
	tmp := []*string{}
	for _, v := range binds {
		if v.RoleID == roleID {
			tmp = append(tmp, &v.NodeID)
		}
	}
	list := []*entities.Node{}
	for _, v := range tmp {
		node, err := nr.nodeRepo.Get(*v)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		list = append(list, node)
	}
	return list, nil
}

func (nr *NodeRoleBindPersistenceUseCase) Update(id string, nodeID string, roleID string) error {
	fmt.Println("update bind", id)
	return nr.nrbRepo.Update(id, nodeID, roleID)
}

func (nr *NodeRoleBindPersistenceUseCase) Delete(id string) error {
	fmt.Println("delete bind", id)
	return nr.nrbRepo.Delete(id)
}

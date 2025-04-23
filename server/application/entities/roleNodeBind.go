package entities

type RoleNodeBind struct {
	ID     string
	NodeID string
	RoleID string
}

// RoleNodeBind is a constructor for initializing a new instance of RoleNodeBind.
func NewRoleNodeBind(id string, nodeID string, roleID string) *RoleNodeBind {
	return &RoleNodeBind{
		ID:     id,
		NodeID: nodeID,
		RoleID: roleID,
	}
}

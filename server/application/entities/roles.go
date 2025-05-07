package entities

import (
	"errors"
	"fmt"
)

// Role represents a role in the system with permissions.
type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

// NewRole creates a new Role entity with parameter validation.
func NewRole(id, name string, permissions ...string) (*Role, error) {
	// Validate ID and Name
	if id == "" || name == "" {
		return nil, fmt.Errorf("ID and Name are required")
	}

	// If no permissions, initialize with an empty list
	if len(permissions) == 0 {
		permissions = []string{}
	} else {
		// Validate each permission
		for _, p := range permissions {
			if len(p) < 3 || len(p) > 50 {
				return nil, fmt.Errorf("Permission '%s' must be between 3 and 50 characters", p)
			}
		}
	}

	// Create and return the Role
	return &Role{
		ID:          id,
		Name:        name,
		Permissions: permissions,
	}, nil
}

// Validate checks the integrity of a Role's data.
func (r *Role) Validate() error {
	if r.ID == "" || r.Name == "" {
		return errors.New("ID and Name are required")
	}
	for _, p := range r.Permissions {
		if len(p) < 3 || len(p) > 50 {
			return fmt.Errorf("permission '%s' must be between 3 and 50 characters", p)
		}
	}
	return nil
}

func (r *Role) HasPerm(str string) bool {
	for _, v := range r.Permissions {
		if v == str {
			return true
		}
	}
	return false
}

//go:build host
// +build host

package handlers

import (
	"fmt"
	"server2/application/entities"
	"server2/application/useCases/handlers/commands"
)

type HandleCommandsUseCase struct {
	ResponseRepo entities.ResponsesReporisory
}

func (r *HandleCommandsUseCase) Execute(command entities.Command) (string, error) {
	if command.IsEncrypted {
		if command.Entry == "add" {
			return commands.Add(command.Id, command.Args, r.ResponseRepo)
		}
	}
	return "", fmt.Errorf("command not found: %s", command.Raw)
}

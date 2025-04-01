//go:build !host
// +build !host

package handlers

import (
	"fmt"
	"server2/application/entities"
	"server2/application/infrastructure"
	"server2/application/useCases/handlers/commands"
)

type HandleCommandsUseCase struct {
	ResponseRepo *infrastructure.ResponsesReporisory
}

func (r *HandleCommandsUseCase) Execute(command entities.Command) (string, error) {
	if command.IsEncrypted {
		switch command.Entry {
		case "block":
			return commands.Block(command.Id, false, command.Args)
		case "unblock":
			return commands.Block(command.Id, true, command.Args)
		case "list":
			return commands.List(command.Id, command.Args)
		case "reload":
			err := commands.ReloadConfig(command.Id)
			return "", err
		case "redirect":
			return commands.AddRedirect(command.Id, command.Args)
		case "unredirect":
			return commands.RemoveRedirect(command.Id, command.Args)
		case "add":
			return commands.Add(command.Id, command.Args, *r.ResponseRepo)
		default:
			return "", fmt.Errorf("comando desconhecido: %s", command.Entry)
		}
	}
	return "", fmt.Errorf("command not found: %s", command.Raw)
}

//go:build !host
// +build !host

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
		if command.Entry == "block" {
			return commands.Block(command.Id, false, command.Args)
		}
		if command.Entry == "unblock" {
			return commands.Block(command.Id, true, command.Args)

		}
		if command.Entry == "list" {
			return commands.List(command.Id, command.Args)

		}
		if command.Entry == "reload" {
			err := commands.ReloadConfig(command.Id)
			return "", err
		}
		if command.Entry == "redirect" {
			return commands.AddRedirect(command.Id, command.Args)

		}
		if command.Entry == "unredirect" {
			return commands.RemoveRedirect(command.Id, command.Args)

		}
		if command.Entry == "add" {
			return commands.Add(command.Id, command.Args, r.ResponseRepo)
		}
	}
	return "", fmt.Errorf("command not found: %s", command.Raw)
}

package commands

import (
	"fmt"
	"server2/application/entities"
	"strings"
)

func Add(id string, args []string, responseRepo entities.ResponsesReporisory) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("wrong Syntax: (%v)\nuse 'id' add <response|connection> data", args)
	}

	if args[0] == "response" {
		// add response(arg0) <id>(arg1) data(rest args concatened)
		data := strings.Join(args[2:], " ")
		responseRepo.Set(args[1], data)
		// memory.ResponseCH <- args[1]
	}

	// fmt.Println(connections)
	return "ok", nil
}

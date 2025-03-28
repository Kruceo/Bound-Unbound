package commands

import (
	"bufio"
	"fmt"
	"os"
	"server2/enviroment"
	"strings"
)

func Block(id string, unblock bool, args []string) (string, error) {

	if len(args) < 1 {
		return "", fmt.Errorf("wrong Syntax: (%v)\nuse 'id' block 'address.net'", args)
	}

	archive, err := os.OpenFile(enviroment.BLOCK_FILEPATH, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	scanner := bufio.NewScanner(archive)
	archiveMap := make(map[string]struct{})

	if err != nil {
		panic(err)
	}

	for scanner.Scan() {
		formatedLine := strings.ReplaceAll(scanner.Text(), "\"", "")
		if formatedLine == "" {
			continue
		}
		address := strings.Split(formatedLine, " ")[1]
		fmt.Println("(" + address + ")")
		archiveMap[address] = struct{}{}
	}

	defer archive.Close()

	if unblock {
		for _, addr := range strings.Split(args[0], ",") {
			delete(archiveMap, addr)
		}
	} else {
		for _, addr := range strings.Split(args[0], ",") {
			archiveMap[addr] = struct{}{}
		}
		// archiveMap[args[0]] = struct{}{}
	}
	archive.Truncate(0)
	for k := range archiveMap {
		fmt.Println("writing: " + k)
		_, err := archive.WriteString(fmt.Sprintf("local-zone: \"%s\" always_nxdomain\n", k))
		if err != nil {
			panic(err)
		}
	}

	// host.AddResponse(id, "ok")
	// fmt.Println(args[0] + " blocked")
	return "ok", nil
}

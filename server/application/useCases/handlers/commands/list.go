package commands

import (
	"bufio"
	"fmt"
	"os"
	"server2/enviroment"
	"server2/utils"
	"strings"
)

func List(id string, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("wrong syntax\nuse list <blocked|redirects>")
	}

	if args[0] == "blocked" {

		archive, err := os.OpenFile(enviroment.BLOCK_FILEPATH, os.O_RDONLY, 0644)
		scanner := bufio.NewScanner(archive)

		if err != nil {
			panic(err)
		}
		response := ""
		for scanner.Scan() {
			formatedLine := strings.ReplaceAll(scanner.Text(), "\"", "")
			if formatedLine == "" {
				continue
			}
			address := strings.Split(formatedLine, " ")[1]
			// address := line[1]
			response += address + ","
		}
		response = strings.TrimSuffix(response, ",")
		fmt.Println(response)
		return response, nil
	}
	if args[0] == "redirects" {

		archive, err := os.OpenFile(enviroment.FORWARD_FILEPATH, os.O_RDONLY, 0644)
		scanner := bufio.NewScanner(archive)

		if err != nil {
			panic(err)
		}
		response := ""
		for scanner.Scan() {
			formatedLine, found := strings.CutPrefix(scanner.Text(), "# @redirect: ")
			if !found {
				continue
			}
			splt := strings.Split(formatedLine, " ")
			from := splt[0]
			recordType := splt[1]
			to := splt[2]
			localZone := splt[3]
			response += fmt.Sprintf("%s %s %s %s,", from, recordType, to, localZone)
		}
		response = strings.TrimSuffix(response, ",")
		fmt.Println(response)
		return response, nil
		// host.AddResponse(id, response)
	}
	if args[0] == "confighash" {

		hash, err := utils.CombinedFileHash([]string{enviroment.BLOCK_FILEPATH, enviroment.FORWARD_FILEPATH})
		if err != nil {
			return "", err
		}

		fmt.Println(hash)
		// err = host.AddResponse(id, hash)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
	}
	return "", nil
}

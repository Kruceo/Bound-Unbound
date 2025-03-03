package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unbound-side-client/host"

	"github.com/gorilla/websocket"
)

func List(conn *websocket.Conn, id string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("wrong syntax\nuse list <blocked|redirects>")
	}
	if args[0] == "blocked" {

		archive, err := os.OpenFile(host.BLOCK_FILEPATH, os.O_RDONLY, 0644)
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
		host.AddResponse(conn, id, response)
	}
	return nil
}

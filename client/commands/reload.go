package commands

import (
	"fmt"
	"os/exec"
	"unbound-side-client/host"

	"github.com/gorilla/websocket"
)

func ReloadConfig(conn *websocket.Conn, id string) error {
	err := exec.Command(host.RELOAD_COMMAND[0], host.RELOAD_COMMAND[1:]...).Run()
	data := "ok"
	if err != nil {
		fmt.Println(err)
		data = "error"
	}
	host.AddResponse(conn, id, data)
	return err
}

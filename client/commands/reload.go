package commands

import (
	"fmt"
	"os/exec"
	"unbound-side-client/host"

	"github.com/gorilla/websocket"
)

func ReloadConfig(conn *websocket.Conn, id string) error {
	err := exec.Command(host.RELOAD_COMMAND[0], host.RELOAD_COMMAND[1:]...).Run()
	if err != nil {
		fmt.Println("host reload command: " + err.Error())
	}

	// hash, err := utils.CombinedFileHash([]string{host.BLOCK_FILEPATH, host.FORWARD_FILEPATH})
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }

	// hashFile, err := os.OpenFile("./blkNfwd.hash", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }
	// defer hashFile.Close()

	// _, err = hashFile.WriteString(hash)

	data := "ok"
	if err != nil {
		fmt.Println(err)
		data = "error"
	}
	host.AddResponse(conn, id, data)
	return err
}

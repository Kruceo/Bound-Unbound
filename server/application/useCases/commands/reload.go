package commands

import (
	"fmt"
	"os/exec"
	"server2/enviroment"
)

func ReloadConfig(id string) error {
	cmd := exec.Command(enviroment.RELOAD_COMMAND[0], enviroment.RELOAD_COMMAND[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("host reload command: " + err.Error())
	}

	fmt.Println(cmd.String(), "\n", string(out))
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

	if err != nil {
		fmt.Println(err)
	}

	return err
}

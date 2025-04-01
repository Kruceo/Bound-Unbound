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

	fmt.Println("executing", cmd.String(), "\n", string(out))

	if err != nil {
		fmt.Println(err)
	}

	return err
}

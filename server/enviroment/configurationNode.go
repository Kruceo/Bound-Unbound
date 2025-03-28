//go:build !host
// +build !host

package enviroment

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"server2/utils"
	"strings"
)

var FORWARD_FILEPATH string
var BLOCK_FILEPATH string
var UNBOUND_CONF_FILEPATH string
var RELOAD_COMMAND []string
var MAIN_SERVER_ADDRESS string

func InitLocals() {
	MAIN_SERVER_ADDRESS = utils.GetEnvOrDefault("MAIN_SERVER_ADDRESS", "127.0.0.1:8080")
	RELOAD_COMMAND = strings.Split(utils.GetEnvOrDefault("UNBOUND_RELOAD_COMMAND", "unbound-control reload"), " ")
	FORWARD_FILEPATH = utils.GetEnvOrDefault("FORWARD_FILEPATH", "/opt/unbound/etc/unbound/forward_records.conf")
	BLOCK_FILEPATH = utils.GetEnvOrDefault("BLOCK_FILEPATH", "/opt/unbound/etc/unbound/block_records.conf")
	UNBOUND_CONF_FILEPATH = utils.GetEnvOrDefault("UNBOUND_CONF_FILEPATH", "/opt/unbound/etc/unbound/unbound.conf")

	file, err := os.OpenFile(FORWARD_FILEPATH, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file, err = os.OpenFile(BLOCK_FILEPATH, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file, err = os.OpenFile(UNBOUND_CONF_FILEPATH, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	if err != nil {
		panic(err)
	}
	configFile := ""
	for scanner.Scan() {
		formatedLine := scanner.Text()
		configFile += formatedLine + "\n"
	}

	r := regexp.MustCompile(`(?s)(# unbound manager configuration)(.+?)(# unbound manager configuration end)`)
	if !r.MatchString(configFile) {
		configFile += "\n\n# unbound manager configuration\n"
		configFile += fmt.Sprintf("server:\n   include: %s\n", FORWARD_FILEPATH)
		configFile += fmt.Sprintf("\n   include: %s\n", BLOCK_FILEPATH)
		configFile += "\n\n# unbound manager configuration end\n"

		file.Truncate(0)
		file.Write([]byte(configFile))

	}

	fmt.Println("="+r.FindString(configFile), r.MatchString(configFile))
}

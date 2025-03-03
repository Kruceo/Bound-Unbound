package host

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unbound-side-client/utils"

	"github.com/joho/godotenv"
)

var FORWARD_FILEPATH string = "/opt/unbound/etc/unbound/forward_records.conf"
var BLOCK_FILEPATH string = "/opt/unbound/etc/unbound/block_records.conf"
var UNBOUND_CONF_FILEPATH string = "/opt/unbound/etc/unbound/unbound.conf"
var RELOAD_COMMAND = []string{"unbound-control", "reload"}
var MAIN_SERVER_ADDRESS = "127.0.0.1:8080"

func init() {
	godotenv.Load(".env")
	FORWARD_FILEPATH = utils.GetEnvOrDefault("FORWARD_FILEPATH", "/opt/unbound/etc/unbound/forward_records.conf")
	BLOCK_FILEPATH = utils.GetEnvOrDefault("BLOCK_FILEPATH", "/opt/unbound/etc/unbound/block_records.conf")
	UNBOUND_CONF_FILEPATH = utils.GetEnvOrDefault("UNBOUND_CONF_FILEPATH", "/opt/unbound/etc/unbound/unbound.conf")
	RELOAD_COMMAND = strings.Split(utils.GetEnvOrDefault("UNBOUND_RELOAD_COMMAND", "unbound-control reload"), " ")
	MAIN_SERVER_ADDRESS = utils.GetEnvOrDefault("MAIN_SERVER_ADDRESS", "127.0.0.1:8080")
	_, err := os.OpenFile(FORWARD_FILEPATH, os.O_RDONLY, 0644)
	if err != nil {
		_, err = os.Create(FORWARD_FILEPATH)
		if err != nil {
			panic(err)
		}
	}

	_, err = os.OpenFile(BLOCK_FILEPATH, os.O_RDONLY, 0644)
	if err != nil {
		_, err = os.Create(BLOCK_FILEPATH)
		if err != nil {
			panic(err)
		}
	}

	f, err := os.OpenFile(UNBOUND_CONF_FILEPATH, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)

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

		f.Truncate(0)
		f.Write([]byte(configFile))

	}

	fmt.Println("="+r.FindString(configFile), r.MatchString(configFile))
}

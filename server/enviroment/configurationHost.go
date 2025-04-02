//go:build host
// +build host

package enviroment

import (
	"fmt"
	"math/rand"
	"server2/utils"
)

var FORWARD_FILEPATH string
var BLOCK_FILEPATH string
var UNBOUND_CONF_FILEPATH string
var RELOAD_COMMAND []string
var MAIN_SERVER_ADDRESS string
var SESSION_SECRET string
var NAME string = "NAMELESS"

func InitLocals() {
	utils.GetEnvOrDefault("SESSION_SECRET", fmt.Sprintf("%x", rand.Int()))
}

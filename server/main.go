package main

import (
	"server2/application/controllers"
	"server2/enviroment"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env")
	enviroment.InitLocals()
}

func main() {
	if controllers.IsHost {
		controllers.RunWebsocketAsHost()
		return
	}
	controllers.RunWebsocketAsNode()
}

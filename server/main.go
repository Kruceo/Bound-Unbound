package main

import (
	"server2/cmd"
	"server2/enviroment"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env")
	enviroment.InitLocals()
}

func main() {
	cmd.Run()
}

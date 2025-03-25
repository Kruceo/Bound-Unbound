package main

import (
	v1 "unbound-mngr-host/api/v1"
	"unbound-mngr-host/enviroment"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env")
	enviroment.InitLocals()
	v1.InitAuth()
}

func main() {
	Run()
}

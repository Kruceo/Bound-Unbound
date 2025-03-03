package host

import "github.com/gorilla/websocket"

func AddResponse(conn *websocket.Conn, id string, data string) error {
	err := conn.WriteMessage(websocket.TextMessage, []byte("add response "+id+" "+data))
	return err
}

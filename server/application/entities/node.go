package entities

import (
	"crypto/cipher"

	"github.com/gorilla/websocket"
)

type Node struct {
	Conn   *websocket.Conn
	Name   string
	Cipher cipher.AEAD
}

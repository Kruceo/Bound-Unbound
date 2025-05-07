package entities

import (
	"crypto/cipher"

	"github.com/gorilla/websocket"
)

type Node struct {
	ID     string
	Conn   *websocket.Conn
	Name   string
	Cipher *cipher.AEAD
}

func NewNode(id string, name string, conn *websocket.Conn, cipher *cipher.AEAD) *Node {
	return &Node{ID: id, Name: name, Conn: conn, Cipher: cipher}
}

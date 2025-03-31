package entities

import (
	"crypto/cipher"
	"fmt"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type Command struct {
	Entry       string
	Args        []string
	Id          string
	IsEncrypted bool
	Raw         string
}

func (command Command) String() string {
	return fmt.Sprintf("%s %s", command.Entry, strings.Join(command.Args, " "))
}

func (command Command) Equal(otherCommand Command) bool {
	return otherCommand.String() == command.String()
}

func (command Command) ArgAsInt(index int) (int32, error) {
	if len(command.Args)-1 < index {
		return 0, fmt.Errorf("this arg not exists: %d", index)
	}
	value, err := strconv.ParseInt(command.Args[index], 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(value), nil
}

func (command Command) ArgAsFloat(index int) (float32, error) {
	if len(command.Args)-1 < index {
		return 0, fmt.Errorf("this arg not exists: %d", index)
	}
	value, err := strconv.ParseFloat(command.Args[index], 32)
	if err != nil {
		return 0, err
	}
	return float32(value), nil
}

type Node struct {
	Conn   *websocket.Conn
	Name   string
	Cipher cipher.AEAD
}

type NodeRepository interface {
	Save(node Node) (string, error)
	Get(id string) *Node
	Delete(id string) error
	IDs() []string
}

type ResponsesReporisory interface {
	Set(id string, data string) error
	WaitForResponse(id string) error
	ReadResponse(id string) (string, error)
	DeleteResponse(id string) error
}

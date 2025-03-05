package commands

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"unbound-side-client/host"
	"unbound-side-client/utils"

	"github.com/gorilla/websocket"
)

func AddRedirect(conn *websocket.Conn, id string, args []string) error {

	// validation
	if len(args) < 3 {
		return fmt.Errorf("wrong Syntax: (%v)\nuse 'id' redirect 'address.net' 'A|AAAA|CNAME' 'other.address.net'", args)
	}

	typeExists, subType := utils.ValidateRecordType(args[1])

	if !strings.Contains(args[0], ".") {
		return fmt.Errorf("the entry is not a domain")
	}

	if !typeExists {
		return fmt.Errorf("the entry is not a record type")
	}
	if subType == "domain" {
		if !strings.Contains(args[2], ".") {
			return fmt.Errorf("the target is not a domain")
		}
	}
	if subType == "ip4" || subType == "ip6" {
		if net.ParseIP(args[2]) == nil {
			return fmt.Errorf("the target is not a valid address")
		}
	}

	if subType == "txt" {
		args[2] = strings.Join(args[2:], " ")
	}

	archive, err := os.OpenFile(host.FORWARD_FILEPATH, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	scanner := bufio.NewScanner(archive)
	type FRedirect struct {
		Type string
		To   string
	}
	archiveMap := make(map[string]FRedirect)

	if err != nil {
		panic(err)
	}

	for scanner.Scan() {
		formatedLine, found := strings.CutPrefix(scanner.Text(), "# @redirect: ")
		if !found {
			continue
		}
		splt := strings.Split(formatedLine, " ")
		from := splt[0]
		recordType := splt[1]
		to := splt[2]

		archiveMap[from] = FRedirect{Type: recordType, To: to} // Adiciona ao mapa sem desperdiçar memória
	}
	// return nil
	defer archive.Close()

	archiveMap[args[0]] = FRedirect{Type: args[1], To: args[2]}

	archive.Truncate(0)
	for k, v := range archiveMap {
		fmt.Println("writing: " + k)
		_, err := archive.WriteString(
			fmt.Sprintf("local-zone: \"%s\" redirect\nlocal-data: \"%s IN %s %s\"\n", k, k, v.Type, v.To) +
				fmt.Sprintf("# @redirect: %s %s %s\n", k, v.Type, v.To))

		if err != nil {
			panic(err)
		}
	}

	// local-zone: "google.com." redirect
	// local-data: "google.com. 3600 IN CNAME facebook.com."

	host.AddResponse(conn, id, "ok")
	// fmt.Println(args[0] + " blocked")
	return nil
}

func RemoveRedirect(conn *websocket.Conn, id string, args []string) error {

	// validation
	if len(args) < 1 {
		return fmt.Errorf("wrong Syntax: (%v)\nuse 'id' unredirect 'address.net'", args)
	}

	if !strings.Contains(args[0], ".") {
		return fmt.Errorf("the entry is not a domain")
	}

	archive, err := os.OpenFile(host.FORWARD_FILEPATH, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	scanner := bufio.NewScanner(archive)
	type FRedirect struct {
		Type string
		To   string
	}
	archiveMap := make(map[string]FRedirect)

	if err != nil {
		panic(err)
	}

	for scanner.Scan() {
		formatedLine, found := strings.CutPrefix(scanner.Text(), "# @redirect: ")
		if !found {
			continue
		}
		splt := strings.Split(formatedLine, " ")
		from := splt[0]
		recordType := splt[1]
		to := splt[2]

		archiveMap[from] = FRedirect{Type: recordType, To: to} // Adiciona ao mapa sem desperdiçar memória
	}
	// return nil
	defer archive.Close()

	delete(archiveMap, args[0])

	archive.Truncate(0)
	for k, v := range archiveMap {
		fmt.Println("writing: " + k)
		_, err := archive.WriteString(
			fmt.Sprintf("local-zone: \"%s\" redirect\nlocal-data: \"%s IN %s %s\"\n", k, k, v.Type, v.To) +
				fmt.Sprintf("# @redirect: %s %s %s\n", k, v.Type, v.To))

		if err != nil {
			panic(err)
		}
	}

	host.AddResponse(conn, id, "ok")
	// fmt.Println(args[0] + " blocked")
	return nil
}

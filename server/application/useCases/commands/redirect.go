package commands

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"server2/enviroment"
	"server2/utils"
	"strings"
)

type FRedirect struct {
	Type      string
	To        string
	LocalOnly bool
}

func AddRedirect(id string, args []string) (string, error) {

	// validation
	if len(args) < 3 {
		return "", fmt.Errorf("wrong Syntax: (%v)\nuse 'id' redirect 'address.net' 'A|AAAA|CNAME' 'other.address.net' 'local-zone|none'", args)
	}

	typeExists, subType := utils.ValidateRecordType(args[1])

	if !strings.Contains(args[0], ".") {
		return "", fmt.Errorf("the entry is not a domain")
	}

	if !typeExists {
		return "", fmt.Errorf("the entry is not a record type")
	}
	if subType == "domain" {
		if !strings.Contains(args[2], ".") {
			return "", fmt.Errorf("the target is not a domain")
		}
	}
	if subType == "ip4" || subType == "ip6" {
		if net.ParseIP(args[2]) == nil {
			return "", fmt.Errorf("the target is not a valid address")
		}
	}

	if subType == "txt" {
		args[2] = strings.Join(args[2:], " ")
	}

	archive, err := os.OpenFile(enviroment.FORWARD_FILEPATH, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	scanner := bufio.NewScanner(archive)

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
		localOnly := false

		if len(splt) >= 4 && splt[3] == "true" {
			localOnly = true
		}
		archiveMap[from] = FRedirect{Type: recordType, To: to, LocalOnly: localOnly} // Adiciona ao mapa sem desperdiçar memória
	}
	// return nil
	defer archive.Close()

	localOnly := false
	if len(args) >= 4 && args[3] == "local-zone" {
		localOnly = true
	}

	archiveMap[args[0]] = FRedirect{Type: args[1], To: args[2], LocalOnly: localOnly}

	archive.Truncate(0)
	for k, v := range archiveMap {
		res := ""
		if v.LocalOnly {
			res += fmt.Sprintf("local-zone: \"%s\" redirect\n", k)
		}
		res += fmt.Sprintf("local-data: \"%s IN %s %s\"\n", k, v.Type, v.To)
		res += fmt.Sprintf("# @redirect: %s %s %s %v\n", k, v.Type, v.To, v.LocalOnly)
		_, err := archive.WriteString(res)

		if err != nil {
			panic(err)
		}
	}

	// local-zone: "google.com." redirect
	// local-data: "google.com. 3600 IN CNAME facebook.com."
	// fmt.Println(args[0] + " blocked")
	return "ok", nil
}

func RemoveRedirect(id string, args []string) (string, error) {

	// validation
	if len(args) < 1 {
		return "", fmt.Errorf("wrong Syntax: (%v)\nuse 'id' unredirect 'address.net'", args)
	}

	if !strings.Contains(args[0], ".") {
		return "", fmt.Errorf("the entry is not a domain")
	}

	archive, err := os.OpenFile(enviroment.FORWARD_FILEPATH, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	scanner := bufio.NewScanner(archive)

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
		localOnly := false

		if len(splt) >= 4 && splt[3] == "true" {
			localOnly = true
		}

		archiveMap[from] = FRedirect{Type: recordType, To: to, LocalOnly: localOnly} // Adiciona ao mapa sem desperdiçar memória
	}
	// return nil
	defer archive.Close()

	delete(archiveMap, args[0])

	archive.Truncate(0)
	for k, v := range archiveMap {
		res := ""
		if v.LocalOnly {
			res += fmt.Sprintf("local-zone: \"%s\" redirect\n", k)
		}
		res += fmt.Sprintf("local-data: \"%s IN %s %s\"\n", k, v.Type, v.To)
		res += fmt.Sprintf("# @redirect: %s %s %s %v\n", k, v.Type, v.To, v.LocalOnly)
		_, err := archive.WriteString(res)

		if err != nil {
			panic(err)
		}
	}

	// host.AddResponse(id, "ok")
	// fmt.Println(args[0] + " blocked")
	return "ok", nil
}

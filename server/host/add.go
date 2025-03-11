package host

import (
	"unbound-mngr-host/memory"
)

func AddResponse(id string, data string) error {
	host, err := memory.GetHost()
	if err != nil {
		return err
	}
	host.Send("_ add response "+id+" "+data, true)
	return nil
}

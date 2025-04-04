# Bound Unbound

The project make the management of multiple Unbound servers easy, concentrating everything in one place.

## Targets

### Important 

- ✅ Easily block domains
- ✅ Easily redirect domains (A,AAAA,CNAME,MX,TXT)
- ✅ Encryption over Host and Unbound node communication
- ✅ Reload Unbound
- ✅ Frontend - Host Auth


### Side Targets
- ✅ Block multiple login requests
- ✅ Store client ip in JWT and compare with client request 
- ✅ Password reset method
- ❌ Multiple users, with custom roles and nodes assigned to each
- ❌ Admin will create "links" to other other person create a guest user
- ❌ Verify if unbound is running after a server reload, if not, restore last config file
- ❌ Modify Unbound configuration remotely
- ❌ Change some host configuration remotely
- ❌ Timed Rules
- ❌ TOTP for two factor authentication

## Build

All tests and builds were run with `go 1.24.1`.

The host and node server share some code together, so the builds are configured using tags. 

### Building Unbound Node binary

```bash
go build -o "bunbound-node"
```

### Building Host Server binary

```bash
go build -tags=host -o "bunbound-host"
```

## Running host

You will need to run the binary with at less one port open (8080).

```bash
go run main.go --host
# ./bound-unbound --host
```

## Running clients (nodes)

Clients will need some config with `.env`.

### Example:

```bash
FORWARD_FILEPATH      = "/opt/unbound/etc/unbound/forward-records.conf"
BLOCK_FILEPATH        = "/opt/unbound/etc/unbound/block-records.conf"
UNBOUND_CONF_FILEPATH = "/opt/unbound/etc/unbound/unbound.conf"
UNBOUND_RELOAD_COMMAND= "unbound-control reload" 
# UNBOUND_RELOAD_COMMAND= "systemctl restart unbound" 
# UNBOUND_RELOAD_COMMAND= "docker exec my-unbound unbound-control reload" # if running with docker
# UNBOUND_RELOAD_COMMAND= "docker restart my-unbound" # if remote control aren't enabled 
NAME                  = "Rafola Telecom"
MAIN_SERVER_ADDRESS   = "127.0.0.1:8080" # your host address and port
```

## WebUI

WebUI will need some config with `.env`.

```bash
# point this to your host
API_PROTOCOL=http
API_ADDRESS=127.0.0.1
API_PORT=8080
```

### Dev

```bash
bun run dev
```

### Build

```bash
bun run build
```

### Running

```bash
bun run prod
```
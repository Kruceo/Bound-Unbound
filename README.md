# Bound Unbound

The project make the management of multiple Unbound servers easy, concentrating everything in one place.

## Targets

### Important 

- ✅ Easily block domains
- ✅ Easily redirect domains (A,AAAA,CNAME,MX,TXT)
- ✅ Encryption over Host - Unbound Client communication
- ✅ Reload Unbound
- ❌ Frontend - Host Auth


### Side Targets

- ❌ Modify Unbound configuration remotely
- ❌ Timed Rules
- ❌ Roles


## Running host

You will need to run the binary with at less one port open (8080); The project uses the same binary to host and client.

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
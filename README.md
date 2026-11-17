# Portscan

A simple concurrent TCP port scanner written in Go.

## Usage

```bash
portscan [flags] <host> <ports>
```

### Examples

```bash
portscan localhost 80
portscan -workers 1 localhost 80,443,25565
portscan -workers 4 localhost 1-1024
portscan -workers 6 localhost -- # for all possible ports
```

## Flags

| Flag | Default | Description |
|---|---:|---|
| `-workers` | `4` | Number of concurrent workers |

## Features

- Concurrent TCP port scanning
- Configurable worker count
- Port ranges, comma-separated ports or all ports
- Connection timeout
- Graceful cancellation with Ctrl+C
- Automatic cancellation on errors

## Build

### Linux
```bash
make build
# or
go build .
```

### Windows
```bash
go build .
```

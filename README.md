# portwatch

Lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with a default scan interval of 30 seconds:

```bash
portwatch start
```

Specify a custom interval and define allowed ports:

```bash
portwatch start --interval 60 --allow 22,80,443
```

When an unexpected port is detected, `portwatch` will log an alert to stdout:

```
[ALERT] 2024/01/15 14:32:01 Unexpected port opened: 4444 (PID: 8821)
[ALERT] 2024/01/15 14:35:10 Port closed: 80
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--interval` | `30` | Scan interval in seconds |
| `--allow` | none | Comma-separated list of allowed ports |
| `--log` | stdout | Path to log file |
| `--quiet` | false | Suppress stdout output |
| `--pid` | none | Path to write the daemon PID file |

### Example with log file

```bash
portwatch start --interval 15 --allow 22,443 --log /var/log/portwatch.log
```

### Stopping the daemon

Send `SIGINT` or `SIGTERM` to gracefully stop the daemon:

```bash
kill $(cat /var/run/portwatch.pid)
```

## Requirements

- Go 1.21+
- Linux or macOS
- Root or `CAP_NET_ADMIN` privileges recommended for full port visibility

## License

MIT — see [LICENSE](LICENSE) for details.

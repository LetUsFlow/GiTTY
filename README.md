# GiTTY
![GiTTY Logo](gitty.ico)

[![Go Report Card](https://goreportcard.com/badge/github.com/LetUsFlow/GiTTY)](https://goreportcard.com/report/github.com/LetUsFlow/GiTTY)

An efficient and user-focused command-line tool developed in Go for simplified SSH management. GiTTY allows you to pre-configure your frequently accessed SSH servers in a structured configuration file. Subsequently, establishing connections to these servers is made seamless by utilizing your default terminal emulator.

Uses [bubbletea](https://github.com/charmbracelet/bubbletea) and [lipgloss](https://github.com/charmbracelet/lipgloss) for terminal UI.

# Configuration (TOML)

GiTTY reads its configuration from `.gitty.toml` in the working directory. Connections are defined as an array of tables under `[[connection]]`.

A minimal example is shown below; copy it into `.gitty.toml` and adjust values.

```toml
[[connection]]
username = "admin"            # required
hostname = "example.com"      # required
comment  = "Test server"      # optional
command  = "echo 'Connected'" # optional
args     = ["-p", "22"]       # optional

[[connection]]
username = "guest"
hostname = "another-server.net"
```

## Building

### Linux
```
go build
```

### Windows
```
windres -o resource.syso resource.rc
go build
```
The `windres` command is optional, but without it, the final executable won't have the GiTTY icon. `windres` is included in `mingw-w64-binutils`.

## License
GPL-3.0

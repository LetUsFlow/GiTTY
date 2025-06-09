# GiTTY
![GiTTY Logo](gitty.ico)

[![Go Report Card](https://goreportcard.com/badge/github.com/LetUsFlow/GiTTY)](https://goreportcard.com/report/github.com/LetUsFlow/GiTTY)

An efficient and user-focused command-line tool developed in Go for simplified SSH management. GiTTY allows you to pre-configure your frequently accessed SSH servers in a structured configuration file. Subsequently, establishing connections to these servers is made seamless by utilizing your default terminal emulator.

Uses [bubbletea](https://github.com/charmbracelet/bubbletea) and [lipgloss](https://github.com/charmbracelet/lipgloss) for terminal UI.

# Example Configuration

GiTTY reads its configuration from a file named `.gitty.json`, which should be located in the same directory where you run the application. This file contains a JSON array of connection profiles.

Below are examples demonstrating the structure of the configuration.

```jsonc
[
  {
    "username": "admin",             // The username for the connection (required)
    "hostname": "example.com",       // The hostname or IP address of the server (required)
    "comment": "Test server",        // An optional comment for your reference
    "command": "echo 'Connected to example.com'", // An optional command to run after successful login
    "args": [                        // An optional array of command-line arguments
      "-p",
      "22"
    ]
  },
  {
    "username": "guest",             // Another connection profile
    "hostname": "another-server.net" // Only the required fields are present
  }
]
```
Important: The comments in this JSON example are for explanation and will cause errors if included in a real configuration file.

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
